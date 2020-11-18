package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func main() {
	parse("testFile.tf")
}

// ServiceConfig creates a service and then can
// be applied to nodes by tags, which are defined
// in the 'apply_to' attribute.
type ServiceConfig struct {
	Type    string   `hcl:"type,label"`
	Name    string   `hcl:"name,label"`
	ApplyTo []string `hcl:"apply_to"`
	Check   string   `hcl:"check"`
}

// Config is the main configurstion struct
// which contains all over blocks, values, etc.
type Config struct {
	Services []ServiceConfig `hcl:"service,block"`
	Runners  []RunnerConfig  `hcl:"runner,block"`
	Remain   hcl.Body        `hcl:",remain"`
}

// RunnerConfig holds information regarding
// a 'type' of runner (ex: Nagios), the 'name' of
// the runner (ex: debian), and the directory
// where the binaries/executables can be found.
type RunnerConfig struct {
	Type       string `hcl:"type,label"`
	Name       string `hcl:"name,label"`
	Directory  string `hcl:"directory"`
	Executable string `hcl:"executable,optional"`
}

func parse(filePath string) {
	var diags hcl.Diagnostics

	log.Printf("loading parser")
	parser := hclparse.NewParser()
	log.Printf("parsing file %q", filePath)
	f, parseDiags := parser.ParseHCLFile(filePath)
	log.Printf("read file")
	if parseDiags.HasErrors() {
		diags = append(diags, parseDiags...)
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,
			parser.Files(),
			78,
			true,
		)
		log.Printf("writing diagnostics")
		wr.WriteDiagnostics(diags)
		log.Fatal(wr)
	}

	var fooInstance Config
	decodeDiags := gohcl.DecodeBody(f.Body, nil, &fooInstance)
	if decodeDiags.HasErrors() {
		diags = append(diags, decodeDiags...)
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,
			parser.Files(),
			78,
			true,
		)
		wr.WriteDiagnostics(diags)
		log.Fatal(decodeDiags.Error())
	}

	fmt.Printf("%#v", fooInstance)
}
