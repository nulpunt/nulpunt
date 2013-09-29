package main

import (
	"log"
)

func main() {
	log.Println("starting nulpunt server application")

	// parse the command-line flags
	initFlags()
	log.Println("flags initialized")

	// init process
	initProcess()
	log.Println("process initialized")

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
