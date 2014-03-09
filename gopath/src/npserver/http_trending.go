package main

import (
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

func getTrendingHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	log.Printf("GetTrendingHandler called\n")
	docs, err := getTrendingDocs(nil, 9)
	if err != nil {
		log.Printf("error getting trending documents %#v\n", err)
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}

	trending := make([]bson.M, len(docs))
	for i, doc := range docs {
		anns := []Annotation{}
		ann, _ := getLatestAnnotation(bson.M{"documentId": doc.ID})
		if ann != nil {
			anns = append(anns, *ann)
		}
		trending[i] = bson.M{"Document": doc, "Annotations": anns}
	}

	err = json.NewEncoder(rw).Encode(bson.M{"trending": trending})
	if err != nil {
		log.Printf("error writing trending data to client: %s\n", err)
		http.Error(rw, "error", http.StatusInternalServerError)
		return
	}
}

func updateTrendingHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	log.Printf("UpdateTrendingHandler called\n")
	rw.WriteHeader(200)
	rw.Write([]byte("Starting\n"))

	cursor := getDocumentsCursor(bson.M{"published": true})

	// for each document
	iterator := cursor.Batch(50).Iter()
	var document Document
	for iterator.Next(&document) {
		rw.Write([]byte(fmt.Sprintf("updating document %#v", document.Title)))

		annotation, err := getLatestAnnotation(bson.M{"documentId": document.ID})
		if err == mgo.ErrNotFound {
			rw.Write([]byte("No annotations found, "))
		}
		score := calculateTrendingScore(&document, annotation)
		err = updateDocumentScore(document.ID, score)
		if err != nil {
			rw.Write([]byte("Error updating document; cannot set score: "))
		}
		rw.Write([]byte(fmt.Sprintf(" score %v\n", document.Score)))
	}
	rw.Write([]byte("Finished\n"))
}
