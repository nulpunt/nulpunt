package main

import (
	"errors"
	"sync"
	"time"
)

var (
	clientSessions     = make(map[string]*ClientSession)
	clientSessionsLock sync.RWMutex
)

var (
	errClientSessionNotFound = errors.New("could not find ClientSession for given key")
)

const clientSessionTimeoutDuration = 1 * time.Minute

// ClientSession defines the session for a given client.
type ClientSession struct {
	sync.Mutex // extends sync.Mutex: CS is to be locked when in use

	key  string    // key for this CS
	ping chan bool // ping chan, true keeps CS alive, false destroys CS

	account *Account // authorized account (when not nil)
}

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

func getClientSession(key string) (*ClientSession, error) {
	// locking
	clientSessionsLock.RLock()
	defer clientSessionsLock.RUnlock()

	// find cs, error on not found
	cs, exists := clientSessions[key]
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

	// timeout sequence
	for {
		select {
		case <-time.After(clientSessionTimeoutDuration):
			return
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
