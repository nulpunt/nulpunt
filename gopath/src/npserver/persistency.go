package main

import (
	"labix.org/v2/mgo"
	"log"
)

// package-wide shared variables pointing to collections in mongodb
var (
	colAccounts *mgo.Collection
)

// initPersistency sets up database connection and initializes col* variables
// it also ensures indexes are existing and will give a fatal error when that fails.
func initPersistency() {
	// dial to localhost mongoDB instance
	mgoConn, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}

	// get "nulpunt" database
	dbNulpunt := mgoConn.DB("nulpunt")

	// get "users" collection
	colAccounts = dbNulpunt.C("accounts")

	// ensure that key "username" is unique for collection "users".
	err = colAccounts.EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	})
	if err != nil {
		log.Fatal(err)
	}
}
