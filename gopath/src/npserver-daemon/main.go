package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {

	initFlags()

	// check permissions
	if os.Getuid() != 0 && os.Geteuid() != 0 {
		fmt.Printf("npserver-daemon should be run as root, have uid=%d and euid=%d\n", os.Getuid(), os.Geteuid())
		os.Exit(1)
	}

	// check start || sto
	if flags.Start == flags.Stop {
		fmt.Println("need --start or --stop flag")
	}

	// setup args for daemon call
	args := []string{
		"--name=npserver",
		"--noconfig",
		"--errlog=/var/log/npserver-daemon.log",
		"--output=/var/log/npserver.log",
		"--pidfile=/run/npserver-daemon.pid",
		"--unsafe",
		"--",
		"/usr/local/bin/npserver",
		fmt.Sprintf("--unix-socket=%s", flags.UnixSocket),
		fmt.Sprintf("--pid-file=%s", flags.PIDFile),
	}

	// append extra args to args
	args = append(args, extraArgs...)

	// start process
	proc, err := os.StartProcess("daemon", args, &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(os.Geteuid()),
				Gid: uint32(os.Getegid()),
			},
		},
	})
	if err != nil {
		fmt.Printf("os/exec returned an error: '%s'\n", err)
		os.Exit(1)
	}

	// wait for daemon to be ready
	_, err = proc.Wait()
	if err != nil {
		fmt.Printf("proc.Wait() failed. %s\n", err)
		os.Exit(1)
	}

	// all good :)
}
