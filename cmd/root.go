package cmd

import (
	"fmt"
	"io"
	"os"

	"keystone/internal/agent"
	"keystone/internal/config"
	"keystone/internal/logger"
	"keystone/internal/tickets"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	version = "v0.1.0"
)

// NewRootCmd creates the root CLI command.
// managerProvider returns an AgentManager given a directory.
// configLoader returns a Config from a path (can be mocked for tests).
// out is optional stdout/stderr writer (useful for testing).
// store is optional ticket store.
func NewRootCmd(
	managerProvider func(string) *agent.AgentManager,
	configLoader func(string) (*config.Config, error),
	out io.Writer,
	store ...*tickets.Store,
) *cobra.Command {

	var ticketStore *tickets.Store
	if len(store) > 0 && store[0] != nil {
		ticketStore = store[0]
	} else {
		ticketStore = tickets.NewStore(tickets.TicketDir)
	}

	var agentsDir string

	cmd := &cobra.Command{
		Use:     "keystone",
		Short:   "Keystone CLI",
		Long:    "Keystone manages AI API access, configuration, and agent operations.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			if verbose && !getJSONFlag(cmd) {
				fmt.Fprintf(cmd.OutOrStdout(), "Verbose mode enabled (version %s)\n", version)
				logger.Info(fmt.Sprintf("Verbose mode enabled (version %s)", version), false)
			}
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := logger.InitDefault(verbose); err != nil {
				PrintError("logger init", fmt.Sprintf("Failed to initialize logger: %v", err), cmd)
				os.Exit(1)
			}

			cfg, err := configLoader(cfgFile)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to load config: %v", err)
				logger.Error(errMsg, getJSONFlag(cmd))
				PrintError("config load", errMsg, cmd)
				os.Exit(1)
			}

			flagVal, _ := cmd.Flags().GetString("agents-dir")
			switch {
			case flagVal != "":
				agentsDir = flagVal
			case cfg.AgentsDir != "":
				agentsDir = cfg.AgentsDir
			default:
				agentsDir = agent.DefaultAgentsDir
			}
		},
	}

	// Output
	if out != nil {
		cmd.SetOut(out)
		cmd.SetErr(out)
	}

	cmd.SetVersionTemplate("{{.Version}}\n")

	// Persistent flags
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	cmd.PersistentFlags().Bool("json", false, "Output results in JSON format")
	cmd.PersistentFlags().String("agents-dir", "", "path to agent YAML definitions (overrides config)")

	// Subcommands
	cmd.AddCommand(
		newTicketCmd(ticketStore),
		newAgentCmd(func() *agent.AgentManager { return managerProvider(agentsDir) }),
		newAgentRegisterCmd(func() *agent.AgentManager { return managerProvider(agentsDir) }, agentsDir),
		newConfigCmd(configLoader),
		newUsageCmd(),
		newWorkflowCmd(func() *agent.AgentManager { return managerProvider(agentsDir) }, ticketStore),
	)

	return cmd
}

// Execute runs the CLI root command using default loader and manager
func Execute() {
	cmd := NewRootCmd(loadManagerWithConfig, config.Load, os.Stdout)
	if err := cmd.Execute(); err != nil {
		logger.Error(err.Error(), false)
		PrintError("keystone", err.Error(), nil)
		os.Exit(1)
	}
}

// loadManagerWithConfig loads the agent manager from a given dir or default.
func loadManagerWithConfig(dir string) *agent.AgentManager {
	if dir == "" {
		dir = agent.DefaultAgentsDir
	}

	mgr := agent.NewManager()
	if err := agent.LoadAgentsFromConfig(mgr, dir); err != nil {
		logger.Warn(fmt.Sprintf("Failed to load agents from %s: %v", dir, err), false)
	}

	agent.LoadDefaultAgent(mgr)
	return mgr
}
