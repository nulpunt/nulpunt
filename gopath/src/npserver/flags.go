package main

import (
	"fmt"
	flagspkg "github.com/jessevdk/go-flags"
)

// options holds the flag settings
// never write to this struct from outside this file
var flags struct {
	Verbose     bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	UnixSocket  bool   `long:"unix-socket" description:"Serve HTTP over unix socket"`
	PIDFilename string `long:"pidfile" description:"PID file for this process" default:"./npserver.pid"`
}

func initFlags() {
	args, err := flagspkg.Parse(&flags)
	if err != nil {
		flagError := err.(*flagspkg.Error)
		if flagError.Type == flagspkg.ErrHelp {
			return
		}
		if flagError.Type == flagspkg.ErrUnknownFlag {
			fmt.Println("Use --help to view all available options.")
			return
		}
		fmt.Printf("Error parsing flags: %s\n", err)
		return
	}

	// check for unexpected arguments
	if len(args) > 0 {
		fmt.Printf("Unknown argument '%s'.\n", args[0])
		return
	}

	//++ do checks
}
