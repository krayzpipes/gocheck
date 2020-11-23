package shell

import (
	"bytes"
	"log"
	"os/exec"
	"syscall"
)


func Check(command string, args []string) (int, string) {
	var out bytes.Buffer
	var exitCode int
	cmd := exec.Command(command, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
				log.Printf("Exit Status: %d", exitCode)
			}
		} else {
			log.Printf(out.String())
			log.Fatal("error when running command: ", err)
		}
	}
	outString := out.String()
	return exitCode, outString
}