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
	r.Body.Close()
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

func sessionPingHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//++ do something useful with the ping?
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

func sessionAuthenticateAccountHandlerFunc(w http.ResponseWriter, r *http.Request) {
	type inDataType struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type outDataType struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		//++ add account details
	}

	outData := &outDataType{}
	defer func() {
		err := json.NewEncoder(w).Encode(outData)
		if err != nil {
			log.Printf("error encoding data for sessionAuthenticateAccount. %s\n", err)
		}
	}()

	inData := &inDataType{}
	err := json.NewDecoder(r.Body).Decode(inData)
	r.Body.Close()
	if err != nil {
		log.Printf("error decoding data for sessionAuthenticateAccount. %s\n", err)
		outData.Error = "Something went wrong."
		return
	}

	cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
	if err != nil {
		log.Printf("error retrieving session.. doesn't exist? %s\n", err)
		outData.Error = "Could not authenticate your account. Session has become invalid. Please refresh the page."
		return
	}
	defer cs.done()

	authenticated, err := cs.authenticateAccount(inData.Username, inData.Password)
	if err != nil {
		log.Printf("could not authenticate. %s\n", err)
		outData.Error = "Could not authenticate your account. Something went wrong."
	}

	if authenticated {
		outData.Success = true
	}
}

func sessionResumeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	type outDataType struct {
		Success  bool   `json:"success"`
		Username string `json:"username,omitempty"`
	}

	outData := &outDataType{}

	defer func() {
		err := json.NewEncoder(w).Encode(outData)
		if err != nil {
			log.Printf("error: could not encode outData to client for sessionResume: %s\n", err)
			http.Error(w, "error", http.StatusInternalServerError)
		}
	}()

	// retrieve session
	cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
	if err != nil {
		log.Printf("error retrieving session.. doesn't exist? %s\n", err)
		return
	}
	defer cs.done()

	// retrieve and check account
	acc := cs.account
	if acc == nil {
		return
	}

	// all done
	outData.Username = acc.Username
	outData.Success = true
}
