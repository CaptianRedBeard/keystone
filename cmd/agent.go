package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"keystone/internal/agent"

	"github.com/spf13/cobra"
)

var (
	providerFlag string
	modelFlag    string
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
		manager := agent.NewManager()
		agent.RegisterSampleAgent(manager)

		fmt.Println("Available agents:")
		for _, a := range manager.List() {
			fmt.Printf(" - %s\n", a.Name)
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

		manager := agent.NewManager()
		agent.RegisterSampleAgent(manager)

		a, err := manager.Get(agentName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🧠 Running agent '%s'\n", agentName)
		fmt.Printf("Provider: %s | Model: %s\n", providerFlag, modelFlag)

		ctx := context.Background()
		response, err := a.Run(ctx, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("💬 Response:")
		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentRunCmd)

	agentRunCmd.Flags().StringVarP(
		&providerFlag,
		"provider", "p",
		"venice",
		"AI provider to use (e.g. venice, openai, anthropic)",
	)
	agentRunCmd.Flags().StringVarP(
		&modelFlag,
		"model", "m",
		"default",
		"Model or endpoint to use for the request",
	)
}
