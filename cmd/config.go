package cmd

import (
	"keystone/internal/config"

	"github.com/spf13/cobra"
)

var showSecrets bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit Keystone configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config.PrintConfig(showSecrets)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "display API keys in output")
}
