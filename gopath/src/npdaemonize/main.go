package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
		fmt.Println("require `--start` or `--stop` flag")
	}
	if len(flags.NP) == 0 {
		fmt.Println("require `--np` <npcmd> flag")
	}

	if flags.Start {
		startDaemon(flags.NP)
	}
	if flags.Stop {
		stopDaemon(flags.NP)
	}

	// all good :)
}

func startDaemon(np string) {
	// setup args for daemon call
	args := []string{
		fmt.Sprintf("--name=%s", np),
		"--noconfig",
		fmt.Sprintf("--errlog=/var/log/npdaemon-%s.log", np),
		fmt.Sprintf("--output=/var/log/%s.log", np),
		fmt.Sprintf("--pidfile=/var/run/%s.pid", np),
		"--unsafe",
		"--",
		fmt.Sprintf("/usr/local/bin/%s", np),
	}

	// append extra args to args
	args = append(args, extraArgs...)

	// start process
	proc, err := os.StartProcess("/usr/bin/daemon", args, &os.ProcAttr{
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

func stopDaemon(np string) {
	pidFileName := fmt.Sprintf("/var/run/%s.pid", np)
	pidFile, err := os.Open(pidFileName)
	if err != nil {
		fmt.Printf("error opening pidFile(%s): %s\n", pidFileName, err)
		killDaemon(np)
		return
	}
	defer pidFile.Close()
	pidBytes, err := ioutil.ReadAll(pidFile)
	if err != nil {
		fmt.Printf("error reading pidFile(%s): %s\n", pidFileName, err)
		killDaemon(np)
		return
	}
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		fmt.Printf("error converting pidFile(%s) contents(%s) to pid number: %s\n", np, string(pidBytes), err)
		killDaemon(np)
		return
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("error finding process with pid(%s, #%d): %s\n", np, pid, err)
		killDaemon(np)
		return
	}
	err = proc.Signal(os.Interrupt)
	if err != nil {
		log.Printf("error sending SIGINT to process(%s, #%d): %s\n", np, pid, err)
		killDaemon(np)
		return
	}
	state, err := proc.Wait()
	if err != nil {
		fmt.Printf("error waiting for process(%s, #%d): %s\n", np, pid, err)
		killDaemon(np)
		return
	}
	if !state.Success() {
		fmt.Printf("%s(%d) was stopped with an error: %s\n", np, pid, err)
		killDaemon(np)
		return
	}

	// all done
	fmt.Printf("%s(%d) stopped successfully\n", np, pid)
}

func killDaemon(np string) {
	fmt.Println("going to kill process (SIGINT) without waiting for it to shut down.")
	// possibly racy.
	// no guarantee that npserver cleaned up before new npserver is being copied/started
	cmd := exec.Command("killall", "-s", "SIGINT", np)
	cmd.Run()
}
