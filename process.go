package main

import (
	"fmt"
	"os"
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

	//++ graceful shutdown on signals (remove pid file)

}
