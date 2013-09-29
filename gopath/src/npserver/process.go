package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
)

// stuff to help manage this process (gracefull shutdown, etc)
func initProcess() {
	// obtain process id
	pid := os.Getpid()

	// create/open pidFile
	pidFile, err := os.OpenFile(flags.PIDFilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_SYNC, 0664)
	if err != nil {
		if err == os.ErrExist {
			fmt.Printf("Could not start nulpunt, pid file exists already (%s)", flags.PIDFilename)
			os.Exit(1)
		}
		fmt.Printf("Could not create pid file. %s\n", err)
		os.Exit(1)
	}

	// write pid number to pidFile
	_, err = pidFile.WriteString(strconv.Itoa(pid))
	if err != nil {
		fmt.Printf("Could not write pid to file. %s\n", err)
		os.Exit(1)
	}

	// close pidFile
	err = pidFile.Close() // ?? correct way to do this?
	if err != nil {
		fmt.Printf("Could not close pid file. %s\n", err)
		os.Exit(1)
	}

	// register signals to channel
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	go func() {
		select {
		case sig := <-sigChan:
			// inform user about received signal
			fmt.Printf("Received signal %v, quitting.\n", sig)

			// remove pid file
			err := os.Remove(flags.PIDFilename)
			if err != nil {
				fmt.Printf("Error cleaning up pid file. %s\n", err)
			}

			// exit with status 0
			os.Exit(0)
		}
	}()

}
