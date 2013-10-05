package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func registerAccountHandlerFunc(w http.ResponseWriter, r *http.Request) {
	type inDataType struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type outDataType struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}

	// defer sending output
	outData := &outDataType{}
	defer func() {
		err := json.NewEncoder(w).Encode(outData)
		if err != nil {
			log.Printf("could not encode data for registerAccoutn request. %s\n", err)
		}
	}()

	inData := &inDataType{}
	err := json.NewDecoder(r.Body).Decode(inData)
	r.Body.Close()
	if err != nil {
		log.Println("could not decode body for registerAccount request. %s\n", err)
		return
	}

	err = registerNewAccount(inData.Username, inData.Email, inData.Password)
	if err != nil {
		if err == errAccountUsernameNotUnique {
			outData.Error = "Username is already taken."
			return
		}
		log.Println("error creating an account: %s\n", err)
		outData.Error = "A server error occured."
		return
	}

	// all done
	outData.Success = true
}
