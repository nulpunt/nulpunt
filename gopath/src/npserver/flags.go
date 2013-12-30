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
	Verbose          bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	UnixSocket       string `long:"unix-socket" description:"Serve HTTP over unix socket"`
	HTTPFiles        string `long:"http-files" description:"location for the http files" default:"./http-files/"`
	HTTPPort         string `long:"http-port" description:"port for HTTP server to listen on" default:"8000"`
	DisableAlphaAuth bool   `long:"disable-alpha-auth" description:"DEPRECATED (was: disable the alpha authentication check)"`
	Environment      string `long:"environment" description:"environment (db/sock) this instance should use"`
}

// initFlags parses the given flags.
// when the user asks for help (-h or --help): the application exists with status 0
// when unexpected flags is given: the application exits with status 1
func initFlags() {
	args, err := goflags.Parse(&flags)
	if err != nil {
		// assert the err to be a flags.Error
		flagError := err.(*goflags.Error)
		if flagError.Type == goflags.ErrHelp {
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

	// check for unexpected arguments
	// when an unexpected argument is given: the application exists with status 1
	if len(args) > 0 {
		fmt.Printf("Unknown argument '%s'.\n", args[0])
		os.Exit(1)
	}

	if flags.DisableAlphaAuth {
		fmt.Printf("DEPRECATED flag --disable-alpha-auth: alpha auth has been stripped from the source code, therefore --disable-alpha-auth is of no use anymore. It will be removed in future release.")
	}
}
