// DB-types.go

// This file contains the definitive mapping between Go structs and Mongo document field names.
// It's called schema in relational parlance.

package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// Account holds information about an account.
// It should not keep data in-memory, but rather write to db directly.
// This type should just be a good wrapper for db read/write functionality
type Account struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Admin    bool          `bson:"admin"`
	Color    string        `bson:"color"` // read out what was created in AccountDetail (at signup)
}

// Account Details for password authentication.
type AccountDetail struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Email    string        `bson:"email"`
	Color    string        `bson:"color"` // Set up here, read out from the db in Account
	Hash     []byte        `bson:"hash"`
	Salt     []byte        `bson:"salt"`
	N        int           `bson:"n"`       // Parameters for the PBKDF2 hashing.
	R        int           `bson:"r"`       // Parameters for the PBKDF2 hashing.
	P        int           `bson:"p"`       // Parameters for the PBKDF2 hashing.
	Remarks  string        `bson:"remarks"` // field for admin remarks about accounts.
}

// type Profile holds the profile properties
type Profile struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Tags     []string      `bson:"tags"` // contains tag.Tag
}

// type Document holds the document properties
type Document struct {
	ID                 bson.ObjectId `bson:"_id"`
	UploaderUsername   string        `bson:"uploaderUsername"`
	UploadFilename     string        `bson:"uploadFilename"`     // original PDF filename.
	UploadGridFilename string        `bson:"uploadGridFilename"` //  Filename into GridFS
	UploadDate         time.Time     `bson:"uploadDate"`
	Language           string        `bson:"language"`
	PageCount          int           `bson:"pageCount"`
	AnalyseState       string        `bson:"analyseState"`

	Title    string   `bson:"title"`
	Summary  string   `bson:"summary"`
	Category string   `bson:"category"`
	Tags     []string `bson:"tags"` // contains tag.Tag

	FOIRequester string    `bson:"foiRequester"`
	FOIARequest  string    `bson:"foiaRequest"`
	OriginalDate time.Time `bson:"orginalDate"`
	Source       string    `bson:"source"`
	Country      string    `bson:"country"`
	Published    bool      `bson:"published"`
}

// type Tag hold the document classification tags
// Note: tags have an ObjectId, these are not for referencing in other collections.
// Just insert the tag-string into other collections where needed.
// Tag: The tag as seen on the site.
type Tag struct {
	ID  bson.ObjectId `bson:"_id"`
	Tag string        `bson:"tag"`
}

type Page struct {
	ID            bson.ObjectId `bson:"_id"`
	DocumentID    bson.ObjectId `bson:"documentId"` // refers to `documents._id`)
	PageNumber    uint          `bson:"pageNumber"` // page number
	Lines         [][]PageChar  `bson:"lines"`
	Text          string        `bson:"text"` // the text in the same order as the lines-attribute, use for search/sharing. Contains ocr-errors
	HighresWidth  int           `bson:"highresWidth"`
	HighresHeight int           `bson:"highresHeight"`
}

type PageChar struct {
	X1 float32 `bson:"x1"` // in percentage relative to the image
	Y1 float32 `bson:"y1"` // in percentage relative to the image
	X2 float32 `bson:"x2"` // in percentage relative to the image
	Y2 float32 `bson:"y2"` // in percentage relative to the image
	C  string  `bson:"c"`
}

// type Annotation hold the document annotations
type Annotation struct {
	ID                bson.ObjectId `bson:"_id"`
	DocumentID        bson.ObjectId `bson:"documentId"` // to which document belong these annotations.
	AnnotatorUsername string        `bson:"annotatorUsername"`
	Color             string        `bson:"color"`
	CreateDate        time.Time     `bson:"createDate"`
	AnnotationText    string        `bson:"annotationText"`
	Locations         []Location    `bson:"locations"`
	Comments          []Comment     `bson:"comments"`
}

type Location struct {
	PageNumber int     `bson:"pageNumber"`
	X1         float32 `bson:"x1"` // in percentage relative to the image
	Y1         float32 `bson:"y1"` // in percentage relative to the image
	X2         float32 `bson:"x2"` // in percentage relative to the image
	Y2         float32 `bson:"y2"` // in percentage relative to the image
}

type Comment struct {
	ID                bson.ObjectId `bson:"_id"`
	CommenterUsername string        `bson:"commenterUsername"`
	Color             string        `bson:"color"`
	CreateDate        time.Time     `bson:"createDate"`
	CommentText       string        `bson:"commentText"`
	Comments          []Comment     `bson:"comments"`
}
