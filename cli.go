package main

import (
	"bytes"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "-alh")
	var out bytes.Buffer
	cmd.Stdout = &out
	log.Printf("Running ls command...")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Command finished with error: %v", err)
	log.Printf("and result is:\n%q", out.String())
}
