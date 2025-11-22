package cmd

import (
	"fmt"

	"keystone/internal/agent"
	"keystone/internal/logger"

	"github.com/spf13/cobra"
)

func newAgentRegisterCmd(managerProvider func() *agent.AgentManager, configDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "register [agentID]",
		Short: "Register a new agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]

			jsonFlag := getJSONFlag(cmd)
			manager := managerProvider()

			// Check if agent already exists
			if _, err := manager.Get(agentID); err == nil {
				msg := fmt.Sprintf("Agent '%s' already exists, skipping registration.", agentID)
				logger.Warn(msg, jsonFlag)
				if jsonFlag {
					Print(map[string]string{"status": "exists", "agentID": agentID}, "", cmd)
				} else {
					Print(nil, msg, cmd)
				}
				return nil
			}

			// Load agent from configuration
			dir := configDir
			if dir == "" {
				dir = agent.DefaultAgentsDir
			}

			lm := agent.NewLifecycleManager(dir, nil)
			if err := lm.LoadAgent(agentID); err != nil {
				errMsg := fmt.Sprintf("Failed to load agent '%s': %v", agentID, err)
				logger.Error(errMsg, jsonFlag)
				PrintError("agent register", errMsg, cmd)
				return err
			}

			msg := fmt.Sprintf("Registered agent '%s' successfully.", agentID)
			logger.Info(msg, jsonFlag)
			if jsonFlag {
				Print(map[string]string{"status": "ok", "agentID": agentID}, "", cmd)
			} else {
				Print(nil, msg, cmd)
			}

			return nil
		},
	}
}
