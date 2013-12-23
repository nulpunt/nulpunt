package nptypes

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// UnknownTODO indicates that a fields' type is not yet determined
type UnknownTODO interface{}

// Account document/entry in the `accounts` collection
type Account struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"` // primary key. e.g. "GeertJohan" in "@GeertJohan", **indexed**
	Email    string        `bson:"email"`    // optional
	Avatar   UnknownTODO   `bson:"avatar"`   // to be decided, link to GridFS file?
	Admin    bool          `bson:"admin"`    // wether user is administrator or ordinary user
	Tags     []string      `bson:"tags"`     // list of tags this user is interested in
}

// Document document/entry in the `documents` collection
type Document struct {
	ID bson.ObjectId `bson:"_id"`

	// upload+analyse information
	UploaderUsername   string    `bson:"uploaderUsername"`   // refers to `accounts.username`
	UploadFilename     string    `bson:"uploadFilename"`     // filename of the original upload
	UploadGridFilename string    `bson:"uploadGridFilename"` // refers to location of the original document in GridFS
	UploadDate         time.Time `bson:"uploadDate"`         // date of upload to nulpunt
	Language           string    `bson:"language"`           // langauge for th document contents, 'nl_NL' format
	PageCount          uint      `bson:"pageCount"`          // number of pages in this document
	AnalyseState       string    `bson:"analyseState"`       // options("uploaded", "started", "completed", "error")

	// document details
	Title        string    `bson:"title"`
	Summary      string    `bson:"summary"`
	Category     string    `bson:"category"`     // "Kamerbrief", "Rapport", ...
	Tags         []string  `bson:"tags"`         // These come from the Tags-table
	FOIRequester string    `bson:"FOIRequester"` // wobber
	FOIARequest  string    `bson:"FOIARequest"`  // wob-verzoek
	OriginalDate time.Time `bson:"originalDate"` // time of publishing by the government agency or date of FOIA-response
	Source       string    `bson:"source"`       // "NL - Binnenlandse zaken", "EN - Foreign affairs", "US - Foreign affairs"
	Country      string    `bson:"country"`      // "NL", "EN"

	// document options
	Published bool `bson:"published"` // true: document is visible for users; false: new or not yet processed document

	// statistics
	Hits uint `bson:"hits"` // number of views for this document
}

// Tag document/entry in the `tags` collection
// Note: tags have an ObjectId, these are not for referencing in other collections.
// Just insert the tag-string into other collections where needed.
type Tag struct {
	ID  bson.ObjectId `bson:"_id"`
	Tag string        `bson:"tag"` // tag
}

// Page document/entry in the `pages` collection
type Page struct {
	ID            bson.ObjectId  `bson:"_id"`
	DocumentID    bson.ObjectId  `bson:"documentId"` // refers to `documents._id`
	PageNumber    uint           `bson:"pageNumber"`
	Lines         []*[]*PageChar `bson:"lines"`
	Text          string         `bson:"text"`          // the text in the same order as the lines-attribute, use for search/sharing. Contains ocr-errors
	HighresWidth  uint           `bson:"highresWidth"`  // the width (in pixels) for the highres(900dpi) render.
	HighresHeight uint           `bson:"highresHeight"` // the height (in pixels) for the highres(900dpi) render.
}

// PageChar is a subtype required for the Page type.
type PageChar struct {
	X1 float32 `bson:"x1"` // offset-left in perecentage relative to the pages' image
	Y1 float32 `bson:"y1"` // offset-top in perecentage relative to the pages' image
	X2 float32 `bson:"x2"` // offset-bottom in perecentage relative to the pages' image
	Y2 float32 `bson:"y2"` // offset-right in perecentage relative to the pages' image
	C  string  `bson:"c"`  // character
}

// Annotation document/entry for the `annotations` collection
type Annotation struct {
	ID                bson.ObjectId        `bson:"_id"`
	DocumentID        bson.ObjectId        `bson:"documentId"`        // refers to `documents._id`
	AnnotatorUsername string               `bson:"annotatorUsername"` // refers to `accounts.username`
	CreateDate        time.Time            `bson:"createDate"`
	Annotation        string               `bson:"annotation"`
	Comments          []Comment            `bson:"comments"`
	Location          []AnnotationLocation `bson:"location"` // in future, there could be multiple sections in a single annotation.
}

// AnnotationLocation is required for the Annotation type.
type AnnotationLocation struct {
	PageNumber uint    `bson:"pageNumber"`
	Y1         float32 `bson:"y1"` // corner-left in percentage relative to the image
	X1         float32 `bson:"x1"` // corner-top in percentage relative to the image
	X2         float32 `bson:"x2"` // corner-bottom in percentage relative to the image
	Y2         float32 `bson:"y2"` // corner-right in percentage relative to the image
}

// Comment is required for the Annotation type.
type Comment struct {
	ID                bson.ObjectId `bson:"_id"`               // needed to do treewalking to get new comments in the right place
	CommenterUsername string        `bson:"commenterUsername"` // refers to `accounts.username`
	CreateDate        time.Time     `bson:"createDate"`
	CommentText       string        `bson:"commentText"`
	Comments          []Comment     `bson:"comments"` // *recursion, disabled for first version??* (JANUARI/FEBRUARI)
}
