package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

//var developer string

var rootCmd = &cobra.Command{
	Use:	"my_flags",
	Short: 	"A brief description of your application",
	Long:	`A longer description of the application`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	//rootCmd.PersistentFlags().StringVar(&developer, "developer", "Unknown Developer!", "Developer name.")
}

func initConfig() {
	//developer, _ := rootCmd.Flags().GetString("developer")
	//if developer != "" {
	//	fmt.Println("developer:", developer)
	//}
}