package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

	if flags.Start {
		startDaemon()
	}
	if flags.Stop {
		stopDaemon()
	}

	// all good :)
}

func startDaemon() {
	// setup args for daemon call
	args := []string{
		"--name=npserver",
		"--noconfig",
		"--errlog=/var/log/npserver-daemon.log",
		"--output=/var/log/npserver.log",
		fmt.Sprintf("--pidfile=%s", flags.PIDFile),
		"--unsafe",
		"--",
		"/usr/local/bin/npserver",
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
}

func stopDaemon() {
	pidFile, err := os.Open(flags.PIDFile)
	if err != nil {
		fmt.Printf("error on opening pidfile: %s", err)
		os.Exit(1)
	}

	pidFileContents, err := ioutil.ReadAll(pidFile)
	pidFile.Close()
	if err != nil {
		fmt.Printf("error reading pidfile contents: %s\n", err)
		os.Exit(1)
	}

	// convert pid string to pid int
	pid, err := strconv.Atoi(string(pidFileContents))
	if err != nil {
		fmt.Printf("error parsing pidfile contents: %s\n", err)
		os.Exit(1)
	}

	// lookup process
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("error finding process with pid %d: %s\n", pid, err)
		os.Exit(1)
	}

	// signal process to stop
	err = proc.Signal(os.Interrupt)
	if err != nil {
		fmt.Printf("error sending interrupt signal to npserver: %s\n", err)
		os.Exit(1)
	}

	// wait until process is done
	state, err := proc.Wait()
	if err != nil {
		fmt.Printf("error waiting for process to stop: %s\n", err)
		os.Exit(1)
	}
	if !state.Exited() || !state.Success() {
		fmt.Printf("npserver process exited badly")
		os.Exit(1)
	}
}
