/*
Package config handles configuration logistics for gocheck.
*/
package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CronType string

// Config is the main configuration built from config
// files.
type Config struct {
	Checks []CheckConfig `hcl:"check,block"`
	Remain hcl.Body      `hcl:",remain"`
}

func (cfg Config) ExtendChecks(other Config) {
	cfg.Checks = append(cfg.Checks, other.Checks...)
}

// CheckConfig maps a check to the devices or applications
// the check should apply to.
type CheckConfig struct {
	Type    string         `hcl:"type,label"`
	Name    string         `hcl:"name,label"`
	ApplyTo []string       `hcl:"apply_to"`
	Cron    string         `hcl:"cron,optional"`
	Exec    CheckExecBlock `hcl:"exec,block"`
}

// CheckExecBlock holds information about the
// executable and arguments needed to perform
// the check.
type CheckExecBlock struct {
	Path string   `hcl:"path"`
	Args []string `hcl:"args"`
}

// ParseConfigFile looks for and parses HCL files used by Gocheck.
func ParseConfigFile(filePath string, cliMode bool) (Config, map[string]*hcl.File, hcl.Diagnostics) {
	var configInstance Config
	var diags hcl.Diagnostics

	parser := hclparse.NewParser()

	f, parseDiags := parser.ParseHCLFile(filePath)

	if parseDiags.HasErrors() {
		diags = append(diags, parseDiags...)
		if cliMode {
			wr := hcl.NewDiagnosticTextWriter(
				os.Stdout,
				parser.Files(),
				78,
				true,
			)
			_ = wr.WriteDiagnostics(diags)
		}
		return configInstance, parser.Files(), diags
	}

	decodeDiags := gohcl.DecodeBody(f.Body, nil, &configInstance)
	if decodeDiags.HasErrors() {
		diags = append(diags, decodeDiags...)
		if cliMode {
			wr := hcl.NewDiagnosticTextWriter(
				os.Stdout,
				parser.Files(),
				78,
				true,
			)
			_ = wr.WriteDiagnostics(diags)
		}
	}
	return configInstance, parser.Files(), diags
}

// GetConfigRootDir Gets the configured root directory for the
// gocheck configs.
func GetConfigRootDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := os.Getenv("GOCHECK_CONFIG_DIR")
	if path == "" {
		path = currentDir
	}
	return path, nil
}

// Recursively walks a directory and looks for gocheck configuration
// files (ending in `.hcl`)
func WalkConfigDirs(root string, ext string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileExt := filepath.Ext(path)
		if fileExt == ext {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// GetParsedConfigFiles gathers configuration and returns a
// configuration object to be used by GoCheck.
func GetParsedConfigFiles(cliMode bool) Config {
	configRootDir, err := GetConfigRootDir()
	if err != nil {
		log.Fatalf("error while looking for configuration directory: %q", err)
	}
	configFiles, err := WalkConfigDirs(configRootDir, ".hcl")
	if err != nil {
		log.Fatalf("error when locating config files: %q", err)
	}

	var rootConfig Config
	for _, configFile := range configFiles {
		anotherConfig, files, diags := ParseConfigFile(configFile, cliMode)
		if 0 < len(diags) {
			var keys []string
			for k := range files {
				keys = append(keys, k)
			}
			keyString := strings.Join(keys, ", ")
			log.Printf("Files not added to running config due to errors: %q", keyString)
		}
		rootConfig.ExtendChecks(anotherConfig)
	}
	return rootConfig
}
