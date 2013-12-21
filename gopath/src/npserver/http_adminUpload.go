package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type uploadedFileMetadata struct {
	UploaderUsername string  `json:"uploaderUsername" bson:"uploaderUsername"`
	Filename         string  `json:"filename" bson:"filename"`
	GridFilename     string  `json:"gridFilename" bson:"gridFilename"` // Consists of: timestamp + randomstring + filename. See database.md GridFS section
	Size             int64   `json:"size" bson:"size"`                 //++ TODO: drop this and use size from gridFS instead?
	Language         *string `json:"language" bson:"language"`
}

func adminUpload(w http.ResponseWriter, r *http.Request) {

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

	// parse multipart form
	err = r.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		log.Printf("error parsing multipart form: %s\n", err)
		return
	}

	// loop over fields and files
	for fieldname, files := range r.MultipartForm.File {
		log.Printf("have fieldname %s\n", fieldname)
		for _, file := range files {
			// generate unique name
			gridFilename := fmt.Sprintf("uploads/%s/%s-%s-%s", acc.Username, strconv.FormatInt(time.Now().Unix(), 10), RandomString(10), file.Filename)
			// metadata instance
			uploadedFile := &uploadedFileMetadata{
				UploaderUsername: acc.Username,
				Filename:         file.Filename,
				GridFilename:     gridFilename,
			}

			// save file
			multipartFile, err := file.Open()
			if err != nil {
				log.Printf("error opening multipart file %s: %s\n", file.Filename, err)
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
			defer multipartFile.Close()
			gridFile, err := gridFS.Create(gridFilename)
			if err != nil {
				log.Printf("error creating gridFile for %s: %s\n", file.Filename, err)
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
			defer gridFile.Close()
			size, err := io.Copy(gridFile, multipartFile)
			if err != nil {
				log.Printf("error copying multipartFile to gridFile for file %s: %s\n", file.Filename, err)
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
			uploadedFile.Size = size

			// save metadata
			err = colUploads.Insert(uploadedFile)
			if err != nil {
				log.Printf("error saving uploadedFile data in colUploads: %s\n", err)
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		}
	}

	w.Write([]byte(`ack`))
}
