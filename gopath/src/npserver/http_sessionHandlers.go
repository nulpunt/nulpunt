package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// sessionInitHandlerFunc creates a new session
func sessionInitHandlerFunc(w http.ResponseWriter, r *http.Request) {
	type outDataType struct {
		SessionKey string `json:"sessionKey"`
	}

	// create a new cs
	newCS := newClientSession()

	// return cs key
	out := &outDataType{
		SessionKey: newCS.key,
	}
	err := json.NewEncoder(w).Encode(out)
	if err != nil {
		log.Printf("Could not encode output data for service sessionInitHandlerFunc. %s\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

// sessionCheckHandlerFunc checks if session is ok
func sessionCheckHandlerFunc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "session ok")
}

// sessionDestroyHandlerFunc destroys a session
func sessionDestroyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	cs.destroy()
	cs.done()
}
