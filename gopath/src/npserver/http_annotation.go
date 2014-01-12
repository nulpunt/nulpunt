package main

import (
	"encoding/json"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"time"
)

// type Annotation struct is defined in annotation.go

func addAnnotationHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("\n\naddAnnotation-request: %v\n", req.URL)

	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("\n\nbody is %s\n", string(body))
		annot := &Annotation{}
		err := json.Unmarshal(body, annot)
		if err != nil {
			log.Printf("\n\nJSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		// Set every other field to things we control.
		annot.ID = bson.NewObjectId()
		annot.AnnotatorUsername = acc.Username
		annot.Color = acc.Color
		annot.CreateDate = time.Now()
		annot.Comments = []Comment{}

		log.Printf("\n\nAnnotation to insert is: %#v\n", *annot)

		err = insertAnnotation(annot)
		if err != nil {
			log.Printf("Error inserting annotation: error %#v\n", err)
			http.Error(rw, "error inserting annotation", http.StatusInternalServerError) // 500
			return
		}

		// marshal and write out.
		j, err := json.Marshal(annot)
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

// AddCommentParams
type AddCommentParams struct {
	//DocumentID bson.ObjectId
	AnnotationID bson.ObjectId
	// Parent bson.ObjectId
	CommentText string
}

// Add a comment to an annotation (non threaded for now)
func addCommentHandler(rw http.ResponseWriter, req *http.Request) {
	log.Printf("\n\naddComment-request: %v\n", req)

	// get session
	cs, err := getClientSession(req.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get account
	acc := cs.account
	if acc == nil {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	switch req.Method {
	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		log.Printf("\n\nbody is %s\n", string(body))
		params := &AddCommentParams{}
		err := json.Unmarshal(body, params)
		if err != nil {
			log.Printf("\n\nJSON unmarshal error %#v\n", err)
			http.Error(rw, "JSON unmarshal error", http.StatusBadRequest) // 400
			return
		}

		comment := Comment{
			ID:                bson.NewObjectId(),
			CommenterUsername: acc.Username,
			Color:             acc.Color,
			CreateDate:        time.Now(),
			CommentText:       params.CommentText,
			Comments:          []Comment{},
		}
		log.Printf("\n\ncomment to add is: %#v\n", comment)

		err = updateAnnotationID(params.AnnotationID, bson.M{"$push": bson.M{"comments": comment}})
		if err != nil {
			log.Printf("Error adding comment to annotation: error %#v\n", err)
			http.Error(rw, "error adding comment to annotation", http.StatusInternalServerError) // 500
			return
		}

		// marshal and write out.
		j, err := json.Marshal(comment)
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
