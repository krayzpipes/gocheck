package cmd

import (
	"github.com/spf13/cobra"
	"gocheck/pkg/config"
	"gocheck/pkg/exec"
	"log"
	"strings"
)

var runCmd = &cobra.Command{
	Use:	"run",
	Short:	"Start a single run of the gocheck agent",
	Long: 	`Use this command to run the services defined in
the configuration of gocheck. The services will be run once and
will ignore any 'apply_to' statements. All configured services
will be run on the machine in which this command is run.`,
	Run:	func(cmd *cobra.Command, args []string) {
		//fmt.Println("run was called")
		fileName, _ := cmd.Flags().GetString("config_file")

		if fileName == "~/.gocheck/default.hcl" {
			log.Printf("No config file identified, will use default config file.")
		}

		configInstance, _, _ := config.ParseConfigFile(fileName, true)

		for _, service := range configInstance.Services {
			check := service.Check
			serviceType := service.Type
			serviceName := service.Name
			checkNameParts := strings.Split(check.Name, ".")

			var executeString string
			for _, checkConfig := range configInstance.Checks {
				if checkConfig.Type == checkNameParts[0] && checkConfig.Name == checkNameParts[1] {
					executeString = checkConfig.Executable
				}
			}
			if executeString == "" {
				log.Fatalf("%v.%v - No valid check found.", serviceType, serviceName)
			}
			exitCode, stdOut := exec.Check(executeString, check.Args)

			if exitCode == 0 {
				log.Printf("%v.%v - HEALTHY: %q", serviceType, serviceName, stdOut)
			} else if exitCode == 1 {
				log.Printf("%v.%v - WARNING: %q", serviceType, serviceName, stdOut)
			} else if exitCode == 2 {
				log.Printf("%v.%v - CRITICAL: %q", serviceType, serviceName, stdOut)
			} else {
				log.Printf("%v.%v - UNKNOWN: %q", serviceType, serviceName, stdOut)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().String("config_file", "~/.gocheck/default.hcl", "The gocheck configuration file to use")
}