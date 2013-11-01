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
}

func (a *Account) getDetails() (*AccountDetails, error) {
	ad := &AccountDetails{}
	err := colAccounts.Find(bson.M{"username": a.Username}).One(ad)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (a *Account) verifyPassword(password string) (bool, error) {
	ad, err := a.getDetails()
	if err != nil {
		return false, err
	}
	return ad.Password == password, nil
}

type AccountDetails struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        // username
	Email    string        // email
	Password string        // password
}

func getAccount(username string) (*Account, error) {
	acc := &Account{}
	err := colAccounts.Find(bson.M{"username": username}).One(acc)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		return nil, err
	}
	return acc, nil
}

func registerNewAccount(username string, email string, password string) error {
	acc := &AccountDetails{
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
