package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"log"
	"os"
)

// Config is the main configuration built from config
// files.
type Config struct {
	Services []ServiceConfig `hcl:"service,block"`
	Checks   []CheckConfig   `hcl:"check,block"`
	Remain   hcl.Body        `hcl:",remain"`
}

// ServiceConfig maps a check to the devices
// the check should apply to.
type ServiceConfig struct {
	Type    string            `hcl:"type,label"`
	Name    string            `hcl:"name,label"`
	ApplyTo []string          `hcl:"apply_to"`
	Check   ServiceCheckBlock `hcl:"check,block"`
}

// ServiceCheckBlock holds information about the
// runner and the arguments passed to the runner
// for a specific check.
type ServiceCheckBlock struct {
	Name string   `hcl:"name"`
	Args []string `hcl:"args"`
}

// CheckConfig holds information on the
// check to be performed
type CheckConfig struct {
	Type       string `hcl:"type,label"`
	Name       string `hcl:"name,label"`
	Executable string `hcl:"executable"`
}


func ParseConfigFile(filePath string) Config {
	var diags hcl.Diagnostics

	parser := hclparse.NewParser()

	f, parseDiags := parser.ParseHCLFile(filePath)

	if parseDiags.HasErrors() {
		diags = append(diags, parseDiags...)
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,
			parser.Files(),
			78,
			true,
		)
		_ = wr.WriteDiagnostics(diags)
		log.Fatal(wr)
	}

	var configInstance Config
	decodeDiags := gohcl.DecodeBody(f.Body, nil, &configInstance)
	if decodeDiags.HasErrors() {
		diags = append(diags, decodeDiags...)
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,
			parser.Files(),
			78,
			true,
		)
		_ = wr.WriteDiagnostics(diags)
		log.Fatal(decodeDiags.Error())
	}

	fmt.Printf("%#v", configInstance)
	return configInstance
}
