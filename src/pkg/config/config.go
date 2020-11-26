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

// Config is the main configuration built from config
// files.
type Config struct {
	Services []ServiceConfig `hcl:"service,block"`
	Checks   []CheckConfig   `hcl:"check,block"`
	Remain   hcl.Body        `hcl:",remain"`
}

func (cfg Config) ExtendServices(other Config) {
	cfg.Services = append(cfg.Services, other.Services...)
}

func (cfg Config) ExtendChecks(other Config) {
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

func WalkConfigDirs(root string, ext string) ([]string, error) {
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
		rootConfig.ExtendServices(anotherConfig)
		rootConfig.ExtendChecks(anotherConfig)
	}
	return rootConfig
}
