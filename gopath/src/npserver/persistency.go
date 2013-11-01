package main

import (
	"labix.org/v2/mgo"
	"log"
)

// package-wide shared variables pointing to collections in mongodb
var (
	colAccounts *mgo.Collection
	gridFS      *mgo.GridFS
)

// initPersistency sets up database connection and initializes col* variables
// it also ensures indexes are existing and will give a fatal error when that fails.
func initPersistency() {
	// dial to localhost mongoDB instance
	mgoConn, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf("fatal error while dialing mgo connection: %s\n", err)
	}

	// get "nulpunt" database
	dbNulpunt := mgoConn.DB("nulpunt")

	// get gridfs
	gridFS = dbNulpunt.GridFS("fs")

	// get "users" collection
	colAccounts = dbNulpunt.C("accounts")

	// ensure that key "username" is unique for collection "users".
	err = colAccounts.EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	})
	if err != nil {
		log.Fatalf("fatal error when ensuring index on accounts.username: %s\n", err)
	}
}
