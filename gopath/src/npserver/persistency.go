package main

import (
	"labix.org/v2/mgo"
	"log"
)

// package-wide shared variables pointing to collections in mongodb
var (
	colAccounts *mgo.Collection
	colProfiles *mgo.Collection
	// colUploads     *mgo.Collection
	colTags        *mgo.Collection
	colDocuments   *mgo.Collection
	colPages       *mgo.Collection
	colAnnotations *mgo.Collection
	colTrending    *mgo.Collection
	gridFS         *mgo.GridFS
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
	// For testing accounts use --environment=  (leave it empty to connect to 'nulpunt').
	database := "nulpunt"
	log.Printf("Environment is: %#v\n", flags.Environment)
	if flags.Environment != "" {
		database += "-" + flags.Environment
	}
	if flags.Verbose == true {
		log.Printf("Connecting to database: %s\n", database)
	}
	dbNulpunt := mgoConn.DB(database)

	// get gridfs
	gridFS = dbNulpunt.GridFS("fs")

	// get "accounts" collection
	colAccounts = dbNulpunt.C("accounts")

	// ensure that key "username" is unique for collection "accounts".
	err = colAccounts.EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	})
	if err != nil {
		log.Fatalf("fatal error when ensuring index on accounts.username: %s\n", err)
	}

	// get "profiles" collection
	colProfiles = dbNulpunt.C("profiles")

	// ensure that key "username" is unique for collection "profiles".
	err = colProfiles.EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	})
	if err != nil {
		log.Fatalf("fatal error when ensuring index on profiles.username: %s\n", err)
	}

	// get "tags" collection
	colTags = dbNulpunt.C("tags")

	// ensure that key "tag" is unique for collection "tags".
	err = colTags.EnsureIndex(mgo.Index{
		Key:    []string{"tag"},
		Unique: true,
	})
	if err != nil {
		log.Fatalf("fatal error when ensuring index on tag.tag: %s\n", err)
	}

	// get "Documents" collection
	colDocuments = dbNulpunt.C("documents")
	err = colDocuments.EnsureIndex(mgo.Index{
		Key: []string{"tags"}, //++ TODO: (published, tags) ?
	})
	if err != nil {
		log.Fatalf("fatal error when ensuting index on documents.tags: %s\n", err)
	}

	// get "Pages" collection
	colPages = dbNulpunt.C("pages")
	err = colPages.EnsureIndex(mgo.Index{
		Key:    []string{"documentId", "pageNumber"},
		Unique: true,
	})
	if err != nil {
		log.Fatalf("fatal error when ensuring index on pages.(documentId, pageNumber): %s\n", err)
	}

	// get "Annotations" collection
	colAnnotations = dbNulpunt.C("annotations")

	// get "trending" collection
	colTrending = dbNulpunt.C("trending")

}
