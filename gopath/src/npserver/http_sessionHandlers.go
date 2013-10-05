package main

import (
	"encoding/json"
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
	type inDataType struct {
		SessionKey string `json:"sessionKey"`
	}

	type outDataType struct {
		Valid bool `json:"valid"`
	}

	inData := &inDataType{}
	err := json.NewDecoder(r.Body).Decode(inData)
	if err != nil {
		log.Printf("error decoding data for sessionCheck. %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&outDataType{
		Valid: isValidClientSession(inData.SessionKey),
	})
	if err != nil {
		log.Printf("error encoding data for sessionCheck. %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
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
