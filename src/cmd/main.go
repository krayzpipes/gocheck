package cmd

import (
	"gocheck/pkg/config"
	"gocheck/pkg/shell"
	"log"
	"strings"
)

func Execute() {
	fileName := "/home/krayzpipes/dev/repos/gocheck/testFile.hcl"

	configInstance := config.ParseConfigFile(fileName)

	for _, service := range configInstance.Services {
		check := service.Check
		checkNameParts := strings.Split(check.Name, ".")

		var executeString string
		for _, checkConfig := range configInstance.Checks {
			if checkConfig.Type == checkNameParts[0] && checkConfig.Name == checkNameParts[1] {
				executeString = checkConfig.Executable
			}
		}
		if executeString == "" {
			log.Fatal("No valid check found.")
		}
		shell.Check(executeString, check.Args)
	}
}