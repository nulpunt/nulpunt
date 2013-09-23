package main

// Account holds information about an account.
// It should not keep data in-memory, but rather write to db directly.
// This type is just a nice wrapper for db functionality
type Account struct {
	username string
}
