package main

import (
	"net/http"
)

func registerAccountHandlerFunc(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	RegisterNewAccount(username, email, password)
}
