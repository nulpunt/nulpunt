package main

import (
	"log"
)

func main() {
	log.Println("starting nulpunt server application")

	setupPersistency()
	log.Println("persistency set up")

	setupHTTPServer()
	log.Println("http server set up")

	log.Println("running..")
	// wait forever
	select {}
}
