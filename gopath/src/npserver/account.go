package main

import (
	"bytes"
	"code.google.com/p/go.crypto/scrypt"
	CryptoRand "crypto/rand"
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
	Username string
}

func (a *Account) getDetails() (*AccountDetail, error) {
	ad := &AccountDetail{}
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
	return ad.ValidatePassword(password), nil
}

type AccountDetail struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        // username
	Email    string        // email
	Hash     []byte
	Salt     []byte
	N, R, P  int    // Parameters for the PBKDF2 hashing.
	Remarks  string // field for admin remarks about accounts.
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
	acc := NewAccountDetail(username, password, email)
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

// Cryptographically strong hash generator.
// Create a new account, salt and hash the password. return it
func NewAccountDetail(username, password, email string) *AccountDetail {
	salt := randBytes(32)
	acct := &AccountDetail{
		ID:       bson.NewObjectId(),
		Username: username,
		Email:    email,
		Salt:     salt,
		N:        16384,
		R:        8,
		P:        1,
	}
	hash := acct.HashPassword(password)
	acct.Hash = hash
	return acct
}

func (a *AccountDetail) HashPassword(password string) []byte {
	hash, err := scrypt.Key([]byte(password), a.Salt, a.N, a.R, a.P, 32)
	if err != nil {
		panic(err)
	}
	return hash
}

func (acct *AccountDetail) ValidatePassword(password string) bool {
	hash := acct.HashPassword(password)
	return bytes.Equal(hash, acct.Hash)
}

func randBytes(length int) (bytes []byte) {
	bytes = make([]byte, length)
	CryptoRand.Read(bytes)
	return
}
