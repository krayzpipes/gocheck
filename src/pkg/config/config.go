package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"log"
	"os"
	"path/filepath"
)

// Config is the main configuration built from config
// files.
type Config struct {
	Services []ServiceConfig `hcl:"service,block"`
	Checks   []CheckConfig   `hcl:"check,block"`
	Remain   hcl.Body        `hcl:",remain"`
}

func (cfg Config) MergeServices(other Config) {
	cfg.Services = append(cfg.Services, other.Services...)
}

func (cfg Config) MergeChecks(other Config) {
	cfg.Checks = append(cfg.Checks, other.Checks...)
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

	//fmt.Printf("%#v", configInstance)
	return configInstance
}

func getConfigRootDir() string {
	currentDir, err := os.Getwd()
	if err == nil {
		log.Fatalf("unable to resolve current working directory: %q", err)
	}
	path := os.Getenv("GOCHECK_CONFIG_DIR")
	if path != "" {
		path = currentDir
	}
	return path
}

func walkConfigDirs(root string, ext string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error{
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

func ParseConfigFiles() Config {
	configRootDir := getConfigRootDir()
	configFiles, err := walkConfigDirs(configRootDir, ".hcl")
	if err != nil {
		log.Fatalf("error when locating config files: %q", err)
	}

	var rootConfig Config
	for _, configFile := range configFiles {
		subConfig := ParseConfigFile(configFile)
		rootConfig.MergeServices(subConfig)
		rootConfig.MergeChecks(subConfig)
	}
	return rootConfig
}
