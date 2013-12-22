package main

import (
	"fmt"
	"os"
	"os/signal"
)

var processEndFuncs = make([]func(), 0)

// stuff to help manage this process (gracefull shutdown, etc)
func initProcess() {

	// register signals to channel
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill) //++ does this really do anything?

	// start a goroutine to wait for and handle a signal
	go func() {
		select {
		case sig := <-sigChan:
			// inform user about received signal
			fmt.Printf("Received %s signal, quitting.\n", sig)

			// call all processEndFuncs
			for _, endFunc := range processEndFuncs {
				endFunc()
			}
			// exit with status 0
			os.Exit(0)
		}
	}()
}
