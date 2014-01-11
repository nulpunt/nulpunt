package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

func downloadOriginalHandlerFunc(w http.ResponseWriter, r *http.Request) {
	type uploadInfoType struct {
		UploadFilename     string `bson:"uploadFilename"`
		UploadGridFilename string `bson:"uploadGridFilename"`
	}

	// get vars from RequestURI
	vars := mux.Vars(r)
	documentIDHex := vars["documentIDHex"]
	filename := vars["filename"]

	// find document
	uploadInfo := &uploadInfoType{}
	err := colDocuments.FindId(bson.ObjectIdHex(documentIDHex)).One(uploadInfo)
	if err != nil {
		if err == mgo.ErrNotFound {
			// document does not exist
			http.NotFound(w, r)
			return
		}
		http.Error(w, "error", http.StatusInternalServerError)
		log.Printf("error retrieving upload info for document %s: %s\n", documentIDHex, err)
		return
	}

	// check filename and redirect on mismatch
	//++ TODO: this works, but doesnt seem to update the URI in my browsers (firefox, chrome).
	//++		I'd think this would be the correct way to do it, although the effect is small and barely visible
	if uploadInfo.UploadFilename != filename {
		http.Redirect(w, r, fmt.Sprintf("/download-original/%s/%s", documentIDHex, uploadInfo.UploadFilename), http.StatusTemporaryRedirect)
		return
	}

	// retrieve uploadedFile from GridFS
	uploadedFile, err := gridFS.Open(uploadInfo.UploadGridFilename)
	if err != nil {
		if err == mgo.ErrNotFound {
			// this IS an error, the original file should be available!
			http.Error(w, "error", http.StatusInternalServerError)
			log.Printf("error! original/uploaded file could not be found for (existing) document %s\n", documentIDHex)
			return
		}
		http.Error(w, "error", http.StatusInternalServerError)
		log.Printf("error when looking up original/uploaded file for (existing) document %s: %s\n", documentIDHex, err)
		return
	}
	defer uploadedFile.Close()

	// set response headers
	w.Header().Add("Content-Type", "application/pdf")
	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, uploadInfo.UploadFilename))

	// write data
	io.Copy(w, uploadedFile)
}
