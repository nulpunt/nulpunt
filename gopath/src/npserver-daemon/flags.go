package main

import (
	"fmt"
	goflags "github.com/jessevdk/go-flags" // rename import to `goflags` (file scope) so we can use `var flags` (package scope)
	"os"
)

// flags holds the flags and arguments given to this process
// never write to this struct from outside this file
// flags is filled with defaults from the tags and the initFlags function
var flags struct {
	Start      bool   `long:"start" description:"start npserver"`
	Stop       bool   `long:"stop" description:"stop npserver"`
	UnixSocket string `long:"unix-socket" description:"unix socket on which the npserver should listen" default:"/var/run/npserver.sock"`
	PIDFile    string `long:"pid-file" description:"pid file for the npserver" default:"/var/run/npserver.pid"`
}

var extraArgs []string

// initFlags parses the given flags.
// when the user asks for help (-h or --help): the application exists with status 0
// when unexpected flags is given: the application exits with status 1
func initFlags() {
	args, err := goflags.Parse(&flags)
	if err != nil {
		// assert the err to be a flags.Error
		flagError := err.(*goflags.Error)
		if flagError.Type == goflags.ErrHelp {
			fmt.Println("npserver-daemon wraps npserver with daemon functionality.")
			// user asked for help on flags.
			// program can exit successfully
			os.Exit(0)
		}
		if flagError.Type == goflags.ErrUnknownFlag {
			fmt.Println("Use --help to view all available options.")
			os.Exit(1)
		}
		fmt.Printf("Error parsing flags: %s\n", err)
		os.Exit(1)
	}
	extraArgs = args
}
