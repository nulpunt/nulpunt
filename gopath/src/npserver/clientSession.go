package main

import (
	"errors"
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
const clientSessionTimeoutDuration = 1 * time.Minute

// ClientSession defines the session for a given client.
type ClientSession struct {
	sync.Mutex // extends sync.Mutex: CS is to be locked when in use

	key  string    // key for this CS
	ping chan bool // ping chan, true keeps CS alive, false destroys CS

	account *Account // authorized account (when not nil)
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
		key:  key,
		ping: make(chan bool),
	}

	// store cs in map
	clientSessions[key] = cs

	// spawn goroutine for ClientSession's life
	go cs.life()

	// all done
	return cs
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

	// lock cs, to be unlocked by cs.done()
	cs.Lock()

	// return cs
	return cs, nil
}

// life keeps track of this sessions lifetime (timeout) and cleans up on destory (cs.ping <- false)
func (cs *ClientSession) life() {
	// end of life story
	defer func() {
		clientSessionsLock.Lock()
		delete(clientSessions, cs.key)
		clientSessionsLock.Unlock()
	}()

	// timeout+ping loop
	for {
		select {
		// when the timeout happens we return and the deferred end-of-life story is started
		case <-time.After(clientSessionTimeoutDuration):
			return

		// when ping value is false we return, and the deferred end-of-life story is started
		// when ping value is true, we just continue to the next loop iteration
		case val := <-cs.ping:
			if !val {
				return
			}
		}
	}
}

// done must be called when the user of ClientSession has no more reads or writes to make.
func (cs *ClientSession) done() {
	cs.Unlock()
}
