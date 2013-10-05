package main

import (
	"labix.org/v2/mgo/bson"
)

// Account holds information about an account.
// It should not keep data in-memory, but rather write to db directly.
// This type should just be a good wrapper for db raed/write functionality
type Account struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        // username
	Email    string        // email
	Password string        // password
}

func RegisterNewAccount(username string, email string, password string) error {
	acc := &Account{
		Username: username,
		Email:    email,
		Password: password,
	}

	// insert into collection
	err := colAccounts.Insert(acc)
	if err != nil {
		return err
	}

	// all done
	return nil
}
