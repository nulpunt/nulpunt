package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func sessionDataBlobSave(w http.ResponseWriter, r *http.Request) {
	type inDataType struct {
		Name string `json:"name"`
		Blob string `json:"blob"`
	}

	cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	acc := cs.account
	if acc == nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	inData := &inDataType{}
	err = json.NewDecoder(r.Body).Decode(inData)
	defer r.Body.Close()
	if err != nil {
		log.Printf("error decoding body for sessionDataSave: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	file, err := gridFS.Create(fmt.Sprintf("/blobs/%s/%s", acc.Username, inData.Name))
	if err != nil {
		log.Printf("error creating GridFile for blob: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	n, err := file.Write([]byte(inData.Blob))
	if err != nil {
		log.Printf("error writing to blob GridFile: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	log.Printf("wrote %d bytes\n", n)

	return
}

func sessionDataBlobLoad(w http.ResponseWriter, r *http.Request) {
	type inDataType struct {
		Name string `json:"name"`
	}

	cs, err := getClientSession(r.Header.Get(headerKeySessionKey))
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer cs.done()

	acc := cs.account
	if acc == nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	inData := &inDataType{}
	err = json.NewDecoder(r.Body).Decode(inData)
	defer r.Body.Close()
	if err != nil {
		log.Printf("error decoding body for sessinoDataBlobLoad: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	file, err := gridFS.Open(fmt.Sprintf("/blobs/%s/%s", acc.Username, inData.Name))
	if err != nil {
		log.Printf("error opening GridFile for blob: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	dataBlob, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading all from GridFile: %s\n ", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(dataBlob)
	if err != nil {
		log.Printf("error writing blob data to client: %s\n", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	return
}
