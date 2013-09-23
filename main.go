package main

import (
	"log"
)

func main() {
	log.Println("starting nulpunt server application")

	// setup connection to mongodb, check indexes, etc.
	initPersistency()
	log.Println("persistency initialized")

	// start a http server
	initHTTPServer()
	log.Println("http server initialized")

	// all seems good
	log.Println("running..")
	// wait forever
	select {}
}
