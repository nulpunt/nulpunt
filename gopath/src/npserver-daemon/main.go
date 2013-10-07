package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Printf("uid: %d\neuid: %d\n", os.Getuid(), os.Geteuid())

	args := []string{
		"--name=npserver",
		"--noconfig",
		"--output=/var/log/npserver.log",
		"--pidfile=/run/npserver-daemon.pid",
		"--unsafe",
		"--",
		"/usr/local/bin/npserver",
		"--unix-socket=/var/run/npserver.sock",
		"--pid-file=/var/run/npserver.pid",
		"--http-files=/srv/nightly.nulpunt.nu",
	}

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
		fmt.Printf("exec returned error: '%s'\n", err)
		os.Exit(1)
	}
	proc.Wait()
}
