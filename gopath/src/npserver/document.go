package main

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

var errDocumentNotUnique = errors.New("Document not unique")
var errPageNotUnique = errors.New("Page not unique")

// newDocument returns a new empty one.
func newDocument() *Document {
	return &Document{
		ID: bson.NewObjectId(),
	}
}

// getDocument gets a single document based upon the specified selection.
// Selection must lead to a unique document. Otherwise, results are undefined
func getDocument(selection interface{}) (*Document, error) {
	doc := &Document{}
	err := colDocuments.Find(selection).One(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// getDocuments gets all documents based upon the specified selection.
// Selection must lead to a unique document. Otherwise, results are undefined
func getDocuments(selection interface{}) ([]Document, error) {
	docs := []Document{}
	err := colDocuments.Find(selection).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func getDocumentsCursor(selection interface{}) *mgo.Query {
	return colDocuments.Find(selection)
}

// insertDocument inserts a new document in the DB. Or updates an existing one.
// Document must have a valid ID, eg from newDocument
func insertDocument(doc *Document) error {
	err := colDocuments.Insert(doc)
	if err != nil {
		mgoErr := err.(*mgo.LastError)
		if mgoErr.Code == 11000 {
			return errDocumentNotUnique
		}
		return err
	}
	// all done
	return nil
}

// upsertDocument inserts a new document in the DB. Or updates an existing one.
// Document must have a valid ID, eg from newDocument
func upsertDocument(docID bson.ObjectId, doc interface{}) error {
	chinfo, err := colDocuments.UpsertId(docID, doc)
	log.Printf("Upsert:change info: %#v\n", chinfo)
	if err != nil {
		mgoErr := err.(*mgo.LastError)
		if mgoErr.Code == 11000 {
			return errDocumentNotUnique
		}
		return err
	}
	// all done
	return nil
}

// removeDocument removes a document from the DB
func removeDocument(docID bson.ObjectId) error {
	err := colDocuments.RemoveId(docID)
	return err
}

//===========================================================

// newPage returns a new empty one.
func newPage() *Page {
	return &Page{
		ID: bson.NewObjectId(),
	}
}

// getPage gets a single document based upon the specified selection.
// Selection must lead to a unique document. Otherwise, results are undefined
func getPage(selection interface{}) (*Page, error) {
	doc := &Page{}
	err := colPages.Find(selection).One(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// getPages gets all documents based upon the specified selection.
// Selection must lead to a unique document. Otherwise, results are undefined
func getPages(selection interface{}) ([]Page, error) {
	docs := []Page{}
	err := colPages.Find(selection).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// insertPage inserts a new document in the DB.
// Page must have a valid ID, eg from newPage
func insertPage(doc *Page) error {
	err := colPages.Insert(doc)
	if err != nil {
		mgoErr := err.(*mgo.LastError)
		if mgoErr.Code == 11000 {
			return errPageNotUnique
		}
		return err
	}
	// all done
	return nil
}

// removePages remove pages from the DB
func removePages(selection interface{}) error {
	err := colPages.Remove(selection)
	// removing a page that does not exist is NOT an error. It's just not there.
	if err == mgo.ErrNotFound {
		return nil
	}
	return err
}
