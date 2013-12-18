package main

import (
	"github.com/GeertJohan/go.leptonica"
	"github.com/GeertJohan/go.tesseract"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

// API
//
// letterbox.Enqueue(document letterbox.Document)
// Accepts a document and places it in the parsing queue.
// All needed data is stored inside a mongo-record.

// letterbox.RunQueue()
// Processess all documents in the queue, asynchronously
// It is assumed that the queue runner lives in a separate process than the np-server.  P
// robably on a different machine.

// Document is a base type in the queue. Both incoming as outgoing queues
type Document struct {
	ID           bson.ObjectId `bson:"_id"`
	DocumentName string
	Language     string
	Image        []byte
	// Add more document field to your tasting.
	Pages      []Page // will be filled after succesful parsing.
	ParseError error  // will be set at errors
}

// Pages are part of documents.
type Page struct {
	PageNumber string // yes, string. it can be anything from the document: "ix", "12", "appendix A-3"
	FullText   string // the unformatted full text for searching
	BoxText    tesseract.BoxText
}

type Queue struct {
	In  *mgo.Collection
	Out *mgo.Collection
}

// Create a new queue on the mongo-db on hostname
func NewQueue(host string) *Queue {
	mgoConn, err := mgo.Dial(host)
	check(err)
	db := mgoConn.DB("nulpunt")
	return &Queue{
		In:  db.C("documentParsingQueueIn"),
		Out: db.C("documentParsingQueueOut"),
	}
}

func (que *Queue) Enqueue(doc Document) error {
	err := que.In.Insert(doc)
	return err
}

// Rename this to Enqueue, add some parameter to taste.
func main() {
	queue := NewQueue("localhost")

	// prepare data
	imageName := "pkiTaskforce.png"
	image, err := ioutil.ReadFile(imageName)
	check(err) // report this error to the user!

	newDoc := Document{
		ID:           bson.NewObjectId(),
		Image:        image,
		DocumentName: "Overheid PKI aanbesteding", // the 'official' name of the document
		Language:     "nld",
		// Set more field to your taste.
	}

	check(queue.Enqueue(newDoc))
	// prentend end of the upload-functionality.
	//return

	//-------------------------------------------------------------------------------------------------------------------

	// Meanwhile at a server far far away:

	//queue := NewQueue("localhost")
	queue.RunQueue()
	return
}

func (que *Queue) RunQueue() {
	query := que.In.Find(nil)
	query.Batch(1) // documents can be big, parsing takes long, don't waste memory.
	iter := query.Iter()

	doc := &Document{}
	for iter.Next(&doc) {
		ProcessDocument(que, doc) // ProcessDoc saves the results in the DB.
	}
	err := iter.Close()
	if err != nil {
		panic(err)
	}
}

// Processes a single document.
func ProcessDocument(que *Queue, doc *Document) {
	if doc.ParseError != nil {
		return // skip previous errors that are stuck in the queue.
	}
	pages, parserr := ParseDocument(doc)
	if parserr != nil {
		// we keep documents with parse errors in the IN-queue.
		// Only valid parsed documents go to the OUT-queue.
		doc.ParseError = parserr
		err := que.In.Update(bson.M{"ID": doc.ID}, doc) // "ID" spelled correctly?
		if err != nil {
			panic(err)
		} // double error, panic
	}

	// parsing went ok.
	doc.Pages = pages

	// store it in mongo.
	// TODO: Atomically
	err := que.Out.Insert(doc)
	if err != nil {
		panic(err)
	}
	err = que.In.RemoveId(doc.ID)
	if err != nil {
		panic(err)
	}
}

func ParseDocument(doc *Document) ([]Page, error) {
	// create new tess instance and point it to the tessdata location.
	tess, err := tesseract.NewTess("/usr/local/share/tessdata", doc.Language)
	if err != nil {
		return nil, err
	}
	defer tess.Close()

	pix, err := leptonica.NewPixReadMem(&doc.Image)
	if err != nil {
		// report this error to the user!
		log.Printf("Error while getting pix from file: %s\n", err)
		return nil, err
	}
	tess.SetImagePix(pix)

	// retrieve boxed text from the tesseract instance
	bt, err := tess.BoxText(0)
	if err != nil {
		log.Printf("Error could not get boxtext from tesseract: %s\n", err)
		return nil, err
	}

	return []Page{Page{
		PageNumber: "1",
		FullText:   tess.Text(),
		BoxText:    *bt,
	}}, nil
}

// Utils

// Check gives panics on error.
// Use this where there is no way to recover from the error other than calling the programmer.
// It's not for user errors, such as selecting wrong files, that they can correct.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
