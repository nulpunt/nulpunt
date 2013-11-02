package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func adminUpload(w http.ResponseWriter, r *http.Request) {

	type uploadedFileType struct {
		Filename   string
		UniqueName string //++ randomstring + filename
		Size       int64
		//++ more metadata
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
			uniqueName := strconv.FormatInt(time.Now().Unix(), 10) + "-" + RandomString(10) + "-" + file.Filename
			// metadata instance
			uploadedFile := &uploadedFileType{
				Filename:   file.Filename,
				UniqueName: uniqueName,
			}

			// save file
			multipartFile, err := file.Open()
			if err != nil {
				log.Printf("error opening multipart file %s: %s\n", file.Filename, err)
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
			defer multipartFile.Close()
			gridFile, err := gridFS.Create(fmt.Sprintf("/uploads/%s/%s", acc.Username, uniqueName))
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
