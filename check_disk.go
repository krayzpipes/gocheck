package main

import (
	"bytes"
	"log"
	"os/exec"
	"syscall"
)

// CheckResult correlates to exit codes
type CheckResult int

const (
	OK       CheckResult = 0
	WARNING              = 1
	CRITICAL             = 2
	UNKNOWN              = 3
)

func main() {
	exitCode, out := check("/usr/lib/nagios/plugins/check_disk", []string{"-w", "100000", "-c", "200000", "-p", "/tmp"})
	if exitCode != 0 {
		log.Printf("FAILED, %q", out)
	} else {
		log.Printf("SUCCESS, %q", out)
	}
}

func check(command string, args []string) (int, string) {
	var out bytes.Buffer
	var exitCode int
	cmd := exec.Command(command, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
				log.Printf("Exit Status: %d", exitCode) // Return exit code, and
			}
		} else {
			log.Printf(out.String())
			log.Printf("error when running command: ", err)
		}
	}
	outString := out.String()
	return exitCode, outString
}
