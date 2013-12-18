package main

import (
	"encoding/json"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

// type Document struct is defined in document.go

type getDocumentParams struct {
	DocID        bson.ObjectId
	AnnotationID bson.ObjectId
	// CommentID bson.ObjectId
}

func getDocumentHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("getDocument-request: %v\n", req)

	// assemble results into a json-object
	result := make(map[string]interface{})

	switch req.Method {
	case "POST":
		// get document, annotation and comment parameters
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("request body is %s\n", string(body))
		params := &getDocumentParams{}
		err := json.Unmarshal(body, params)
		log.Printf("Params is: %#v\n", params)
		if err != nil {
			log.Printf("JSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		if params.DocID == "" {
			log.Printf("DocID is empty.\n")
			http.Error(rw, "DocID is empty", http.StatusBadRequest) // 400
			return
		}

		// get document
		doc, err := getDocument(bson.M{"_id": params.DocID})
		if err != nil {
			log.Printf("DocID not found: error %#v\n", err)
			http.Error(rw, "DocID not found", http.StatusNotFound) // 404
			return
		}
		result["document"] = doc

		// get optional annotation, error if it is specified but not there.
		if params.AnnotationID != "" {
			annotation, err := getAnnotation(bson.M{"_id": params.AnnotationID})
			if err != nil {
				log.Printf("AnnotationID not found: error %#v\n", err)
				http.Error(rw, "AnnotationID not found", http.StatusNotFound) // 404
				return
			}
			result["annotation"] = annotation
		}

		// marshal and write out.
		j, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error marshalling results: error %#v\n", err)
			http.Error(rw, "Marshaling error", http.StatusInternalServerError) // 500
			return
		}
		rw.WriteHeader(200)
		rw.Write(j)
		return

	default: // request.Method
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}

func getDocumentListHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("getDocument-request: %v\n", req)

	switch req.Method {
	case "POST": // Use POST as that's the easiest to encode json parameters
		body, _ := ioutil.ReadAll(req.Body)
		params := &getDocumentParams{}
		err := json.Unmarshal(body, params)
		if err != nil {
			log.Printf("JSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		docs, err := getDocuments(nil)
		if err != nil {
			log.Printf("getDocuments error %#v\n", err)
			http.Error(rw, "getDocuments error", http.StatusInternalServerError) // 400
			return
		}

		j, err := json.Marshal(docs)
		if err != nil {
			log.Printf("Error marshalling results: error %#v\n", err)
			http.Error(rw, "Marshaling error", http.StatusInternalServerError) // 500
			return
		}
		rw.WriteHeader(200)
		rw.Write(j)
		return

	default:
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}

func insertDocumentHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("\n\ninsertDocument-request: %v\n", req)

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("\n\nbody is %s\n", string(body))
		doc := &Document{}
		err := json.Unmarshal(body, doc)
		if err != nil {
			log.Printf("\n\nJSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		log.Printf("\n\nDocument to insert is: %#v\n", *doc)

		if doc.ID == "" {
			doc.ID = bson.NewObjectId()
			log.Printf("\n\nCreating new ObjectId: %v\n", doc.ID)
		}
		err = insertDocument(doc)
		if err != nil {
			log.Printf("Error inserting  document: error %#v\n", err)
			http.Error(rw, "error inserting document", http.StatusInternalServerError) // 500
			return
		}

		//This is a HACK
		// Add page-record
		page := newPage()
		page.DocumentID = doc.ID
		page.PageNr = 1
		page.Text = "Hallo"
		// page.Lines = [][]CharObject{ [ { ...
		err = insertPage(page)
		if err != nil {
			log.Printf("Error inserting page: error %#v\n", err)
			http.Error(rw, "error inserting page", http.StatusInternalServerError) // 500
			return
		}

		rw.WriteHeader(200)
		rw.Write([]byte(`OK, inserted`))
		return
	default:
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}

func updateDocumentHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("\n\nupdateDocument-request: %v\n", req)

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("\n\nbody is %s\n", string(body))
		doc := &Document{}
		err := json.Unmarshal(body, doc)
		if err != nil {
			log.Printf("\n\nJSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		log.Printf("\n\nDocument to update is: %#v\n", *doc)
		err = upsertDocument(doc)
		if err != nil {
			log.Printf("Error inserting/updating  document: error %#v\n", err)
			http.Error(rw, "error inserting/updating document", http.StatusInternalServerError) // 500
			return
		}

		rw.WriteHeader(200)
		rw.Write([]byte(`OK, updated`))
		return
	default:
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}
