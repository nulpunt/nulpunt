package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// during alpha, we use basic auth for alpha access
// basic auth is chosen because it does not interfere with any other techniques used by this project
// the alphaUsers map holds all alpha users and their sha512 hashed passwords.
// the password is to be salted with the alphaSalt (e.g. <password>+alphaSalt)
// feel free to add yourself or other project members to the map
var alphaSalt = "che9ohlu0xie9yaW6aeg"
var alphaUsers = map[string]string{
	"GeertJohan": "8e764667137fc73b16a3a3cd43a6a9314c8c2214215306f563068bd91b7f5de06da19815b10e3b6719c7f9cd4499b26e13738fadba0b244c4055f1d7af0100b8",
}

func alphaCheckBasicAuth(r *http.Request) bool {
	// retrieve auth header
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 {
		// no or invalid auth given
		return false
	}
	if authHeader[:6] != "Basic " {
		// invalid auth type
		return false
	}

	// decode auth data
	authData, err := base64.StdEncoding.DecodeString(authHeader[6:])
	if err != nil {
		return false
	}
	authDataSlice := strings.SplitN(string(authData), ":", 2)
	if len(authDataSlice) != 2 {
		return false
	}
	givenUsername := authDataSlice[0]
	givenPassword := authDataSlice[1]

	// retrieved hashed password from alphaUsers map
	correctPasswordHashed, userExists := alphaUsers[givenUsername]
	if !userExists {
		// user does not exist
		return false
	}

	// hash retrieved password, format as hex string
	passwordHasher := sha512.New()
	io.WriteString(passwordHasher, givenPassword)
	io.WriteString(passwordHasher, alphaSalt)
	givenPasswordHashed := fmt.Sprintf("%x", passwordHasher.Sum(nil))
	// password matches?
	if givenPasswordHashed == correctPasswordHashed {
		// yay! correct auth!
		return true
	}

	return false
}
