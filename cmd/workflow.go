package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"keystone/internal/agent"
	"keystone/internal/logger"
	"keystone/internal/tickets"
	"keystone/internal/workflow"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// newWorkflowCmd creates the "workflow" CLI group and its subcommands
func newWorkflowCmd(managerProvider func() *agent.AgentManager, store *tickets.Store) *cobra.Command {
	wfCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage and run agent workflows",
		Long:  "Run one or more registered agents in a defined workflow sequence.",
	}

	wfCmd.AddCommand(newWorkflowRunCmd(managerProvider, store))
	wfCmd.PersistentFlags().Bool("json", false, "Output results in JSON format")
	wfCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging for workflow steps")

	return wfCmd
}

// newWorkflowRunCmd runs a workflow from a YAML file by workflow ID
func newWorkflowRunCmd(managerProvider func() *agent.AgentManager, store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "run [workflow_id]",
		Short: "Run a workflow by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wfID := args[0]

			jsonFlag, _ := cmd.Flags().GetBool("json")
			verboseFlag, _ := cmd.Flags().GetBool("verbose")

			manager := managerProvider()
			engine := workflow.NewEngine(manager, verboseFlag)

			// Load workflow YAML from workflows/<workflow_id>.yaml
			yamlFile := fmt.Sprintf("workflows/%s.yaml", wfID)
			data, err := os.ReadFile(yamlFile)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to read workflow file '%s': %v", yamlFile, err), jsonFlag)
				return fmt.Errorf("failed to read workflow file: %w", err)
			}

			var wf workflow.Workflow
			if err := yaml.Unmarshal(data, &wf); err != nil {
				logger.Error(fmt.Sprintf("Failed to parse workflow YAML '%s': %v", yamlFile, err), jsonFlag)
				return fmt.Errorf("failed to parse workflow YAML: %w", err)
			}

			// Create a new ticket for this workflow
			ticket := tickets.NewTicket(tickets.NewID("cli", "workflow", wfID), "default", nil)

			// Run the workflow
			results, err := engine.Run(context.Background(), wf, ticket)
			if err != nil {
				logger.Error(fmt.Sprintf("Workflow '%s' run failed: %v", wfID, err), jsonFlag)
				return fmt.Errorf("workflow run failed: %w", err)
			}

			// Log each step output based on verbose flag
			for i, r := range results {
				if r.Error != nil {
					if verboseFlag {
						logger.Error(fmt.Sprintf("Step %d - Agent %s failed: %v", i, r.AgentID, r.Error), jsonFlag)
					}
				} else if verboseFlag {
					logger.Info(fmt.Sprintf("Step %d - Agent %s output:\n%s", i, r.AgentID, r.Output), jsonFlag)
				}
			}

			// Prepare output
			output := map[string]interface{}{
				"results":        results,
				"ticket_context": ticket.Context,
			}

			if jsonFlag {
				enc, _ := json.MarshalIndent(output, "", "  ")
				fmt.Println(string(enc))
			} else {
				fmt.Printf("✅ Workflow '%s' executed successfully\n", wfID)
			}

			return nil
		},
	}
}
