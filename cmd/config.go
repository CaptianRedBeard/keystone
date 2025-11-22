package cmd

import (
	"keystone/internal/config"

	"github.com/spf13/cobra"
)

var showSecrets bool

// newConfigCmd returns a CLI command to view/edit configuration.
// You can inject a custom load function for testing.
func newConfigCmd(loadFn func(string) (*config.Config, error)) *cobra.Command {
	// Use default loader if nil
	if loadFn == nil {
		loadFn = config.Load
	}

	cmd := &cobra.Command{
		Use:   "config",
		Short: "View or edit Keystone configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := loadFn(cfgFile)
			if err != nil {
				PrintError("config load", err.Error(), cmd)
				return
			}

			// TODO: mask secrets if showSecrets is false
			conf := cfg

			if getJSONFlag(cmd) {
				Print(conf, "", cmd)
			} else {
				Print(nil, "Current Keystone configuration", cmd)
				Print(conf, "", cmd)
			}
		},
	}

	cmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "display API keys in output")
	return cmd
}
