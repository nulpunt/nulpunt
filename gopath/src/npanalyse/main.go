package main

import (
	"log"
)

func main() {
	log.Println("starting nulpunt analyse application")

	// parse the command-line flags
	initFlags()
	log.Println("flags initialized")

	// init process
	initProcess()
	log.Println("process initialized")

	// setup connection to mongodb, check indexes, etc.
	initPersistency()
	log.Println("persistency initialized")

	// start analysers
	initAnalysers(flags.NumWorkers)
	log.Println("Analysers initialized")

	// all seems good, inform the user
	log.Println("npanalyse is running..")

	// wait forever
	select {}
}
