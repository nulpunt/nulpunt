package main

import (
	"log"
	"net/http"
)

// setupHTTPServer sets up the http.FileServer and other http services.
func setupHTTPServer() {
	// serve static files from http-files on root
	http.Handle("/", http.FileServer(http.Dir("./http-files/")))

	// run http server in goroutine
	go func() {
		port := "8000"
		log.Printf("starting http server on port %s\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
