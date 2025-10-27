package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	version = "v0.1.0"
)

var rootCmd = &cobra.Command{
	Use:   "keystone",
	Short: "Keystone provides a common entry point for AI agent orchestration",
	Long:  `Keystone manages AI API access, configuration, and agent operations.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Verbose mode enabled")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Keystone operational!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default is $HOME/.keystone/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

// GetConfigPath returns the --config flag value for main.go
func GetConfigPath() string {
	return cfgFile
}
