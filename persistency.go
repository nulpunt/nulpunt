package main

import (
	"labix.org/v2/mgo"
	"log"
)

var (
	colUsers *mgo.Collection
)

func setupPersistency() {
	mgoConn, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}

	dbNulpunt := mgoConn.DB("nulpunt")

	colUsers = dbNulpunt.C("users")
	err := colUsers.EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	})
	if err != nil {
		log.Fatal(err)
	}
}
