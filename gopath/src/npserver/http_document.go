package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

// type Document struct is defined in document.go

type DocumentParams struct {
	DocID        bson.ObjectId
	AnnotationID bson.ObjectId
	// CommentID bson.ObjectId
}

// Get a single document, specified by DocID,
// Get the Annotation, if specified.
func getDocumentHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("getDocument-request: %v\n", req)

	// assemble results into a json-object
	result := make(map[string]interface{})

	switch req.Method {
	case "POST":
		// get document, annotation and comment parameters
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("request body is %s\n", string(body))
		params := &DocumentParams{}
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

		// Be paranoid and limit annotation to the Document they belong to.
		selector := bson.M{"documentId": params.DocID}
		if params.AnnotationID != "" {
			selector["_id"] = params.AnnotationID
		} // or leave it undefined for all annotations of DocID
		annotations, err := getAnnotations(selector)
		if err != nil {
			log.Printf("AnnotationID not found: error %#v\n", err)
			http.Error(rw, "AnnotationID not found", http.StatusNotFound) // 404
			return
		}
		result["annotations"] = annotations

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

// get all documents with certain limits.
// For lazy loading place a start-at at next call
func getDocumentsHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("getDocument-request: %v\n", req)

	// assemble results into a json-object
	result := make(map[string]interface{})

	switch req.Method {
	case "POST":
		// get document, annotation and comment parameters
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("request body is %s\n", string(body))
		params := &DocumentParams{}
		err := json.Unmarshal(body, params)
		log.Printf("Params is: %#v\n", params)
		if err != nil {
			log.Printf("JSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		// get document
		// UGLY HACK: get them all.
		docs, err := getDocuments(nil)
		if err != nil {
			log.Printf("GetDocuments error %#v\n", err)
			http.Error(rw, "GetDocuments error", http.StatusNotFound) // 404
			return
		}
		result["documents"] = docs

		// get optional annotation, error if it is specified but not there.

		// bs := bson.M{}
		// if params.AnnotationID != "" {
		// 	bs = bson.M{"_id": params.AnnotationID}
		// } else {
		// 	bs = bson.M{"DocID": params.DocID}
		// }
		// annotations, err := getAnnotations(bs)
		// if err != nil {
		// 	log.Printf("AnnotationID not found: error %#v\n", err)
		// 	http.Error(rw, "AnnotationID not found", http.StatusNotFound) // 404
		// 	return
		// }
		// result["annotations"] = annotations

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
		params := &DocumentParams{}
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
		// get session
		cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
		if err != nil {
			http.Error(rw, "error", http.StatusInternalServerError)
			return
		}
		defer cs.done()

		// get account
		acc := cs.account
		if acc == nil || acc.Admin == false {
			http.Error(rw, "forbidden", http.StatusForbidden)
			return
		}

		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("\n\nbody is %s\n", string(body))
		doc := &Document{}
		err = json.Unmarshal(body, doc)
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
		page.PageNumber = 1
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
	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil || acc.Admin == false {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

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

func deleteDocumentHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("\n\ndelete-Document-request: %v\n", req)

	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil || acc.Admin == false {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}
	// TODO: test for Admin-user flag.

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("\n\nbody is %s\n", string(body))
		params := &DocumentParams{}
		err := json.Unmarshal(body, params)
		if err != nil {
			log.Printf("\n\nJSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		// Delete Annotation-records with DocID
		err = removeAnnotations(bson.M{"documentId": params.DocID})
		if err != nil {
			log.Printf("Error deleting annotation on document: error %#v\n", err)
			http.Error(rw, "error deleting annotation on document", http.StatusInternalServerError) // 500
			return
		}

		// Delete Page-records with DocID
		err = removePages(bson.M{"documentId": params.DocID})
		if err != nil {
			log.Printf("Error deleting pages of document: error %#v\n", err)
			http.Error(rw, "error deleting pages of document", http.StatusInternalServerError) // 500
			return
		}

		// Delete the document record.
		log.Printf("\n\nDocumentID to delete is: %#v\n", params.DocID)
		err = removeDocument(params.DocID)
		if err != nil {
			log.Printf("Error deleting document: error %#v\n", err)
			http.Error(rw, "error deleting document", http.StatusInternalServerError) // 500
			return
		}

		rw.WriteHeader(200)
		rw.Write([]byte(`OK, deleted`))
		return
	default:
		http.Error(rw, "error", http.StatusMethodNotAllowed) // 405
	}
}

func pageImageHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	urlVars := mux.Vars(r)
	documentIDHex := urlVars["documentIDHex"]
	if !bson.IsObjectIdHex(documentIDHex) {
		http.NotFound(w, r)
		return
	}
	documentID := bson.ObjectIdHex(documentIDHex)
	pageNumberString := urlVars["pageNumber"]
	pageNumber, err := strconv.ParseUint(pageNumberString, 10, 32)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fileName := fmt.Sprintf("docviewer-pages/%s-%d.png", documentID.Hex(), pageNumber)
	file, err := gridFS.Open(fileName)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		log.Printf("error looking up files in gridFS (%s): %s\n", fileName, err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/png")
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("error writing png file (%s) to http client: %s\n", fileName, err)
		return
	}
	// all done :)
}
