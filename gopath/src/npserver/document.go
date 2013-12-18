package main

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

var errDocumentNotUnique = errors.New("Document not unique")
var errPageNotUnique = errors.New("Page not unique")

// type Document holds the document properties
type Document struct {
	ID           bson.ObjectId `bson:"_id"`
	Original     string        // original PDF file to download the whole thing. Filename into GridFS
	Published    bool
	UploadDate   time.Time
	Uploader     string
	Title        string
	Summary      string
	Source       string
	Category     string
	Tags         []string // contains tag.Tag
	OriginalDate time.Time
}

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
func upsertDocument(doc *Document) error {
	chinfo, err := colDocuments.UpsertId(doc.ID, doc)
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
//func removeTag(doc *Document) error {
//	err := colDocument.RemoveId(doc.ID)
//	return err
//}

//===========================================================

// type Page holds the page properties
type Page struct {
	ID         bson.ObjectId `bson:"_id"`
	DocumentID bson.ObjectId
	PageNr     int
	Text       string
	Lines      [][]CharObject
}

type CharObject struct {
	X1   int
	Y1   int
	X2   int
	Y2   int
	Char string
}

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

// removePage removes a page from the DB
//func removePage(page *Page) error {
//	err := colPages.RemoveId(Page.ID)
//	return err
//}
