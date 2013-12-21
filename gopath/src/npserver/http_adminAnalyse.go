package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func adminGetRawUploads(w http.ResponseWriter, r *http.Request) {
	type outDataType struct {
		Files []*uploadedFileMetadata `json:"files"`
	}

	// get session
	cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	// get accuont
	acc := cs.account
	if acc == nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// new outData instance
	outData := &outDataType{
		Files: make([]*uploadedFileMetadata, 0),
	}

	// get uploads that are not analyzed or listed to analyze yet
	err = colUploads.Find(nil).All(&outData.Files)
	if err != nil {
		log.Printf("error retrieving uploaded files from uploads collections: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// send data to client
	err = json.NewEncoder(w).Encode(outData)
	if err != nil {
		log.Printf("error sending uploaded files to client: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// all done
}
