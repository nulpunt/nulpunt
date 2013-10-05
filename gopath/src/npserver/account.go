package main

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var errAccountUsernameNotUnique = errors.New("account username is not unique")

// Account holds information about an account.
// It should not keep data in-memory, but rather write to db directly.
// This type should just be a good wrapper for db raed/write functionality
type Account struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        // username
	Email    string        // email
	Password string        // password
}

func registerNewAccount(username string, email string, password string) error {
	acc := &Account{
		ID:       bson.NewObjectId(),
		Username: username,
		Email:    email,
		Password: password,
	}

	// insert into collection
	err := colAccounts.Insert(acc)
	if err != nil {
		mgoErr := err.(*mgo.LastError)
		if mgoErr.Code == 11000 {
			return errAccountUsernameNotUnique
		}
		return err
	}

	// all done
	return nil
}
