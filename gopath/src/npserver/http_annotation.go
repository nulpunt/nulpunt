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

	body, _ := ioutil.ReadAll(req.Body)
	log.Printf("\n\nbody is %s\n", string(body))
	annot := &Annotation{}
	err = json.Unmarshal(body, annot)
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

	// Normalize the coordinates: X1,Y1 at top left, X2,Y2 as bottom right.
	for i, coord := range annot.Locations {
		coord.X1, coord.X2 = min(coord.X1, coord.X2), max(coord.X1, coord.X2)
		coord.Y1, coord.Y2 = min(coord.Y1, coord.Y2), max(coord.Y1, coord.Y2)
		annot.Locations[i] = coord // set it, as range gives a copy.
	}

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
}

// AddCommentParams holds comment parameters
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

	body, _ := ioutil.ReadAll(req.Body)
	log.Printf("\n\nbody is %s\n", string(body))
	params := &AddCommentParams{}
	err = json.Unmarshal(body, params)
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
}

// Because package Math doesn't have Max and Min for float32...
func min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}
