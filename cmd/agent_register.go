package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"keystone/internal/agent"
	"keystone/internal/logger"

	"github.com/spf13/cobra"
)

var (
	agentID          string
	agentName        string
	agentDescription string
	agentProvider    string
	agentModel       string
	agentMemory      string
	promptTemplate   string
	parametersJSON   string
	logging          bool
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new custom agent",
	Long:  "Dynamically create a new agent, save its config, and register it with Keystone.",
	Run: func(cmd *cobra.Command, args []string) {
		if agentID == "" || agentName == "" {
			fmt.Fprintln(os.Stderr, "Error: --id and --name are required")
			os.Exit(1)
		}

		// Parse parameters JSON
		var params map[string]string
		if parametersJSON != "" {
			if err := json.Unmarshal([]byte(parametersJSON), &params); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse --parameters JSON: %v\n", err)
				os.Exit(1)
			}
		}

		cliCfg := agent.AgentConfig{
			ID:             agentID,
			Name:           agentName,
			Description:    agentDescription,
			Provider:       agentProvider,
			Model:          agentModel,
			Memory:         agentMemory,
			PromptTemplate: promptTemplate,
			Parameters:     params,
			Logging:        logging,
		}

		lm := agent.NewLifecycleManager(GetConfigPath())

		if err := lm.SaveOrMergeConfig(cliCfg); err != nil {
			outputErrorJSON("register", agentID, fmt.Sprintf("Failed to save config: %v", err))
			os.Exit(1)
		}

		if err := lm.LoadAgent(agentID); err != nil {
			outputErrorJSON("register", agentID, fmt.Sprintf("Failed to load agent: %v", err))
			os.Exit(1)
		}

		message := fmt.Sprintf("Agent '%s' registered successfully", agentName)
		if jsonOutput {
			out := map[string]string{
				"status":  "ok",
				"agentID": agentID,
				"message": message,
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf("✅ %s\n", message)
		}
		logger.Log(agentID, fmt.Sprintf("Agent '%s' registered (provider=%s, model=%s)", agentName, agentProvider, agentModel))
	},
}

func init() {
	agentCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVar(&agentID, "id", "", "Unique ID for the agent (required)")
	registerCmd.Flags().StringVar(&agentName, "name", "", "Display name for the agent (required)")
	registerCmd.Flags().StringVar(&agentDescription, "description", "", "Description of the agent")
	registerCmd.Flags().StringVar(&agentProvider, "provider", "venice", "Provider for the agent")
	registerCmd.Flags().StringVar(&agentModel, "model", "default", "Model or endpoint for the agent")
	registerCmd.Flags().StringVar(&agentMemory, "memory", "", "Optional memory/session key for the agent")
	registerCmd.Flags().StringVar(&promptTemplate, "prompt_template", "", "Optional prompt template for the agent")
	registerCmd.Flags().StringVar(&parametersJSON, "parameters", "", "Optional JSON string of parameters for the agent")
	registerCmd.Flags().BoolVar(&logging, "logging", false, "Enable logging for the agent")
}

// Helper for JSON error output
func outputErrorJSON(command, agentID, msg string) {
	if jsonOutput {
		out := map[string]string{
			"status":  "error",
			"agentID": agentID,
			"message": msg,
		}
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Fprintf(os.Stderr, "Error [%s]: %s\n", command, msg)
	}
}
