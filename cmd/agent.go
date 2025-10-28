package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"keystone/internal/agent"
	"keystone/internal/logger"

	"github.com/spf13/cobra"
)

var (
	providerFlag      string
	modelFlag         string
	cliPromptTemplate string
	cliParametersJSON string
	jsonOutput        bool
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage and run Keystone agents",
	Long:  "List, configure, and run Keystone agents.",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available agents",
	Run: func(cmd *cobra.Command, args []string) {
		manager := loadManagerWithConfig()
		agents := manager.List()

		if jsonOutput {
			out := make([]map[string]string, 0, len(agents))
			for _, a := range agents {
				out = append(out, map[string]string{
					"id":          a.ID(),
					"name":        a.Name(),
					"description": a.Description(),
				})
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Println("Available agents:")
		for _, a := range agents {
			fmt.Printf(" - %s: %s\n", a.Name(), a.Description())
		}
	},
}

var agentRunCmd = &cobra.Command{
	Use:   "run [agentName] [input]",
	Short: "Run an agent with an input string",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		agentName := args[0]
		input := strings.Join(args[1:], " ")

		manager := loadManagerWithConfig()
		a, err := manager.Get(agentName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		modelToUse := a.DefaultModel()
		if modelFlag != "" && modelFlag != "default" {
			modelToUse = modelFlag
		}

		provider := a.Provider()
		if providerFlag != "" && providerFlag != "venice" {
			fmt.Println("⚠️ Only Venice provider is currently supported for overrides")
		}

		finalParameters := make(map[string]string)
		for k, v := range a.Parameters() {
			finalParameters[k] = v
		}

		if cliParametersJSON != "" {
			var cliParams map[string]string
			if err := json.Unmarshal([]byte(cliParametersJSON), &cliParams); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse --parameters JSON: %v\n", err)
				os.Exit(1)
			}
			for k, v := range cliParams {
				finalParameters[k] = v
			}
		}

		finalTemplate := a.PromptTemplate()
		if cliPromptTemplate != "" {
			finalTemplate = cliPromptTemplate
		}

		finalInput := input
		if finalTemplate != "" {
			for k, v := range finalParameters {
				finalTemplate = strings.ReplaceAll(finalTemplate, fmt.Sprintf("{{%s}}", k), v)
			}
			finalInput = fmt.Sprintf("%s\n%s", finalTemplate, input)
		}

		ctx := context.Background()
		response, err := provider.GenerateResponse(ctx, finalInput, modelToUse)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			out := map[string]interface{}{
				"agentID":    a.ID(),
				"name":       a.Name(),
				"input":      finalInput,
				"response":   response,
				"parameters": finalParameters,
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Running agent '%s'\nProvider: %T | Model: %s\nResponse:\n%s\n", a.Name(), provider, modelToUse, response)
	},
}

// loadManagerWithConfig loads all agents from the config directory (global or default).
func loadManagerWithConfig() *agent.AgentManager {
	configDir := GetConfigPath()
	if configDir == "" {
		configDir = "./internal/agent/config"
	}

	manager := agent.NewManager()
	if err := agent.LoadAgentsFromConfig(manager, configDir); err != nil {
		logMsg := fmt.Sprintf("Failed to load agents: %v", err)
		logger.Log("system", logMsg)
		fmt.Fprintln(os.Stderr, logMsg)
		os.Exit(1)
	}
	return manager
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentListCmd, agentRunCmd)

	// Global flag for JSON output
	agentCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")

	agentRunCmd.Flags().StringVarP(&providerFlag, "provider", "p", "venice", "AI provider to use")
	agentRunCmd.Flags().StringVarP(&modelFlag, "model", "m", "default", "Model or endpoint to use")
	agentRunCmd.Flags().StringVar(&cliPromptTemplate, "prompt_template", "", "Override prompt template")
	agentRunCmd.Flags().StringVar(&cliParametersJSON, "parameters", "", "Override parameters JSON")
}
