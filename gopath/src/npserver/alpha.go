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
	"Argure":     "cc5be974f66c51babfc1e800847303adf9cb092342dce874483347ce04e9a122a409a9876080f56e3322a2176b490fc72c5e7b87e9fde45b718d09d26d9006dc",
	"nihlaeth":   "efa115c51855f495b31beb31904acc67b223e0d0466cb4f2f42ad7d8b07c43ef25e7ba66ed6b48f97c40ca625d163dede7f2700507602f583874a3da7eabb665",
	"frank88":    "fc647f062eabc707c9bb266ff2a15efab2006e2dee267e2d753da2c5f7e4ce26bc70da916402323012769f5f716accb03420c7952acdd961cf81bedf4365a411",
	"renee":      "1eb07335daaba6d7b86ec7de45e32532eb80e6edfdad5e6a759e57db86d68007c92655873e5bd4bc70fa1b48a16c7173af500c39617e3c027f280dd9ceae1e47",
	"younes":     "996a9b02da5088409c96d609fe5db9e5129efbb335beda06c35e4980814b59a98d4241a4fddf45a2eb93465f520d5ba2cdfe54ed75955bcb43775a88a46db2f2",
	"guido":       "cf283559a5731f10206c8e395bc6d38d1ecf2c5f6e0cc748af9e0793a190a1b2cb1f27baf39b3707961d11cbc9f881297bc307de5f1c493806d1af1e5c4afe88",

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
