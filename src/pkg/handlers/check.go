// Package worker manages logistics with creating checks, running
// checks, scheduling checks, and communicating check results.
package handlers

// Manager to spin up goroutines and listen on channels for responses
// Manager can tell goroutine to stop, this way the program can add
//	more goroutines and not have to reload/restart.

// Set timer/ticker
// Run check
// Set state
// Report results
// Check for instructions

import (
	"fmt"
	"github.com/krayzpipes/cronticker/cronticker"
	"gocheck/pkg/config"
	"gocheck/pkg/exec"
	"log"
	"time"
)

const (
	KILL     = "kill"
	STOPNORM = "stopnorm"
	STOPERR  = "stoperr"
)

// CheckHandler performs checks on a schedule. Listens for
// instructions on a channel from the WorkerManager and returns
// results by the result channel. It can also update the manager on when it
// is shutting down by using the notify channel.
func CheckHandler(check CheckMeta, control <-chan string, result chan<- CheckResult, notify chan<- string) {
	ticker, err := cronticker.NewTicker(check.Cron())
	if err != nil {
		log.Printf("error when creating now cron ticker: %v, cron schedule: '%s'", err, check.Cron)
		notify <- STOPERR
		close(notify)
		return
	}
	defer ticker.Stop()

	for {
		select {
		// If the manager sends instructions, do them
		case c := <-control:
			if c == KILL {
				ticker.Stop()
				notify <- STOPNORM
				close(notify)
				return
			}

		// If the ticker has 'ticked', run the check
		case t := <-ticker.C:
			log.Printf("running check at: %v", t)
			exitCode, outString, errString := exec.Check(check.Executable(), check.Args())
			r := CheckResult{ExitCode: exitCode, StdOut: outString, Error: errString}
			result <- r
			// TODO handle state file?

		default:
			time.Sleep(500 * time.Millisecond)
		}

	}

}

// CheckMeta interface defines the interface expected by the
// CheckHandler function.
type CheckMeta interface {
	Id() string
	Executable() string
	Args() []string
	Cron() string
}

// CheckInfo satisfies the CheckMeta interface which will be used
// to run checks by the CheckHandler.
type CheckInfo struct {
	id         string
	executable string
	arguments  []string
	cron       string
}

func (c *CheckInfo) Id() string {
	return c.id
}

func (c *CheckInfo) Executable() string {
	return c.executable
}

func (c *CheckInfo) Args() []string {
	return c.arguments
}

func (c *CheckInfo) Cron() string {
	return c.cron
}

// NewCheckInfo creates a CheckInfo struct from a
// config.CheckConfig struct.
func NewCheckInfo(c config.CheckConfig) CheckInfo {
	var checkInfo CheckInfo
	checkInfo.id = fmt.Sprintf("%s:%s", c.Type, c.Name)
	checkInfo.executable = c.Exec.Path
	checkInfo.arguments = c.Exec.Args
	checkInfo.cron = c.Cron
	return checkInfo
}
