package main

import (
	"log"
)

func main() {
	log.Println("starting nulpunt server application")

	// setup connection to mongodb, check indexes, etc.
	setupPersistency()
	log.Println("persistency set up")

	// start a http server
	setupHTTPServer()
	log.Println("http server set up")

	// all seems good
	log.Println("running..")
	// wait forever
	select {}
}
