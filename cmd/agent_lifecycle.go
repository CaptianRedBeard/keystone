package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"keystone/internal/agent"
	"keystone/internal/logger"

	"github.com/spf13/cobra"
)

var lifecycleCmd = &cobra.Command{
	Use:   "lifecycle",
	Short: "Manage agent lifecycle (load, unload, reload)",
	Long:  "Commands to dynamically load, unload, or reload agents from disk without restarting Keystone.",
}

var loadCmd = &cobra.Command{
	Use:   "load [agentID]",
	Short: "Load an agent from disk",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleLifecycle("load", args[0], func(lm *agent.LifecycleManager, id string) error {
			return lm.LoadAgent(id)
		})
	},
}

var unloadCmd = &cobra.Command{
	Use:   "unload [agentID]",
	Short: "Unload an agent from memory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleLifecycle("unload", args[0], func(lm *agent.LifecycleManager, id string) error {
			return lm.UnloadAgent(id)
		})
	},
}

var reloadCmd = &cobra.Command{
	Use:   "reload [agentID]",
	Short: "Reload an agent from disk",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleLifecycle("reload", args[0], func(lm *agent.LifecycleManager, id string) error {
			return lm.ReloadAgent(id)
		})
	},
}

func init() {
	agentCmd.AddCommand(lifecycleCmd)
	lifecycleCmd.AddCommand(loadCmd, unloadCmd, reloadCmd)
}

// Unified lifecycle handler
func handleLifecycle(action, agentID string, fn func(*agent.LifecycleManager, string) error) {
	lm := agent.NewLifecycleManager(GetConfigPath())
	if err := fn(lm, agentID); err != nil {
		if jsonOutput {
			out := map[string]string{
				"status":  "error",
				"agentID": agentID,
				"action":  action,
				"message": err.Error(),
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Failed to %s agent '%s': %v\n", action, agentID, err)
		}
		os.Exit(1)
		return
	}

	msg := fmt.Sprintf("Agent '%s' %sed successfully", agentID, action)
	if jsonOutput {
		out := map[string]string{
			"status":  "ok",
			"agentID": agentID,
			"action":  action,
			"message": msg,
		}
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("✅ %s\n", msg)
	}
	logger.Log(agentID, fmt.Sprintf("Agent %s via CLI", action))
}
