package main

// Account holds information about an account.
// It should not keep data in-memory, but rather write to db directly.
// This type should just be a good wrapper for db raed/write functionality
type Account struct {
	username string
}
