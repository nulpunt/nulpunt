package main

import (
	"labix.org/v2/mgo"
	"log"
)

var (
	colUsers *mgo.Collection
)

func setupPersistency() {
	// dial to localhost mongoDB instance
	mgoConn, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}

	// get "nulpunt" database
	dbNulpunt := mgoConn.DB("nulpunt")

	// get "users" collection
	colUsers = dbNulpunt.C("users")

	// ensure that key "username" is unique for collection "users".
	err := colUsers.EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	})
	if err != nil {
		log.Fatal(err)
	}
}
