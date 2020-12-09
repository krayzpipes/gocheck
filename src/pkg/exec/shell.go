package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)


func Check(command string, args []string) (int, string, string) {
	var out bytes.Buffer
	var exitCode int
	var errString string
	cmd := exec.Command(command, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		} else {
			errString = fmt.Sprintf("error when calling command: %s", err)
		}
	}
	outString := out.String()
	return exitCode, outString, errString
}