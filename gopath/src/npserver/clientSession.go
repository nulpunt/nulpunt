package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

// locked map of clientSessions
var (
	clientSessions     = make(map[string]*ClientSession)
	clientSessionsLock sync.RWMutex
)

// internally used errors
var (
	errClientSessionNotFound = errors.New("could not find ClientSession for given key")
)

// configuration constants
const clientSessionTimeoutDuration = 7 * time.Minute

// ClientSession defines the session for a given client.
type ClientSession struct {
	sync.Mutex // extends sync.Mutex: CS is to be locked when in use

	key       string    // key for this CS
	destroyCh chan bool // destroy chan, receive from this chan destroys session
	pingCh    chan bool // ping chan, receive from this chan keeps session alive

	account *Account // authorized account (when nil: no auth)
}

// newClientSession creates a new ClientSession
// a new unique key is generated for this ClientSession
// ClientSession lifecycle is also started by this function
func newClientSession() *ClientSession {
	// locking
	clientSessionsLock.Lock()
	defer clientSessionsLock.Unlock()

	// find unique key
	var key string
	for {
		key = RandomString(30)
		if _, exists := clientSessions[key]; !exists {
			break
		}
	}

	// create new ClientSession instance
	cs := &ClientSession{
		key:       key,
		destroyCh: make(chan bool),
		pingCh:    make(chan bool),
	}

	// store cs in map
	clientSessions[key] = cs

	// spawn goroutine for ClientSession's life
	go cs.life()

	// all done
	return cs
}

func isValidClientSession(key string) bool {
	// locking
	// not using defered lock to avoid long-locking the clientSessions map when
	// 		multiple goroutines require a client (cs.Lock in this function will hang)
	clientSessionsLock.RLock()

	// find cs, error on not found
	cs, exists := clientSessions[key]
	clientSessionsLock.RUnlock()
	if !exists {
		return false
	}

	if !cs.ping() {
		return false
	}

	// all ok
	return true
}

// getClientSession returns a client session by given key (if it can be found)
// when an error occurs or client session is not found, nil+error are returned
// the returned ClientSession is locked, and can only be used by one
func getClientSession(key string) (*ClientSession, error) {
	// locking
	// not using defered lock to avoid long-locking the clientSessions map when
	// 		multiple goroutines require a client (cs.Lock in this function will hang)
	clientSessionsLock.RLock()

	// find cs, error on not found
	cs, exists := clientSessions[key]
	clientSessionsLock.RUnlock()
	if !exists {
		return nil, errClientSessionNotFound
	}

	if !cs.ping() {
		return nil, errClientSessionNotFound
	}

	// lock cs, to be unlocked by cs.done()
	cs.Lock()

	// return cs
	return cs, nil
}

// life keeps track of this sessions lifetime (timeout) and cleans up on destory (cs.destroy <- false)
func (cs *ClientSession) life() {
	// end of life story
	defer func() {
		// close channels (destroy() and ping() will now return false)
		close(cs.destroyCh)
		close(cs.pingCh)

		clientSessionsLock.Lock()
		delete(clientSessions, cs.key)
		clientSessionsLock.Unlock()
	}()

	// timeout+destroy+ping loop
	for {
		select {
		// when the timeout happens we return and the deferred end-of-life story is started
		case <-time.After(clientSessionTimeoutDuration):
			return

		case cs.destroyCh <- true:
			return

		case cs.pingCh <- true:
			continue
		}

	}
}

// ping checks the ClientSession for life
// when false is returned, the ClientSession is not alive
func (cs *ClientSession) ping() bool {
	_, ok := <-cs.pingCh
	return ok
}

// destroy tries to destroy the ClientSession
// when false is returned, it was already destroyed
func (cs *ClientSession) destroy() bool {
	_, ok := <-cs.destroyCh
	return ok
}

// done must be called when the user of ClientSession has no more reads or writes to make.
func (cs *ClientSession) done() {
	cs.Unlock()
}

// Authenticate returns the account record if successful, nil otherwise
// We need the Admin flag and the Color value in the user interface
func (cs *ClientSession) authenticateAccount(username string, password string) (*Account, error) {
	log.Printf("authenticateAccount got: %#v, %#v\n", username, password)
	acc, err := getAccount(username)
	if err != nil {
		return nil, err
	}
	log.Printf("authenticateAccount got account: %#v\n", acc)
	if acc == nil {
		return nil, nil
	}
	valid, err := acc.verifyPassword(password)
	if err != nil {
		return nil, err
	}
	if valid {
		cs.account = acc
		log.Printf("authenticateAccount returns %#v\n", acc)
		return acc, nil
	}
	return nil, nil
}
