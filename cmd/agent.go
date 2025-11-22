package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"keystone/internal/agent"
	"keystone/internal/tickets"

	"github.com/spf13/cobra"
)

// newAgentCmd creates the "agent" command and subcommands
func newAgentCmd(managerProvider func() *agent.AgentManager) *cobra.Command {
	var (
		providerFlag      string
		modelFlag         string
		cliPromptTemplate string
		cliParametersJSON string
		ticketFlag        string
		verboseFlag       bool
	)

	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage and run Keystone agents",
	}

	// -------- list command --------
	agentCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available agents",
		Run: func(cmd *cobra.Command, args []string) {
			manager := managerProvider()
			agents := manager.List()
			out := make([]map[string]string, 0, len(agents))
			for _, a := range agents {
				out = append(out, map[string]string{
					"id":          a.ID(),
					"name":        a.Name(),
					"description": a.Description(),
				})
			}
			printOrJSON(out, fmt.Sprintf("Available agents listed (%d)", len(out)), cmd)
		},
	})

	// -------- run command --------
	runCmd := &cobra.Command{
		Use:   "run [agentName] [input]",
		Short: "Run an agent with an input string (optionally bound to a ticket)",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			input := strings.Join(args[1:], " ")

			mgr := managerProvider()
			a, err := mgr.Get(agentName)
			if err != nil {
				PrintError("agent run", fmt.Sprintf("Agent '%s' not found", agentName), cmd)
				return nil
			}

			ticket, store, err := loadOrCreateTicket(ticketFlag)
			if err != nil {
				return fmt.Errorf("ticket error: %w", err)
			}

			finalParams, err := mergeCLIParams(a.Parameters(), cliParametersJSON, cmd)
			if err != nil {
				return err
			}

			finalInput, usedTemplate := applyPromptTemplate(a.PromptTemplate(), cliPromptTemplate, finalParams, input)

			ctx := context.Background()
			resp, err := a.Handle(ctx, finalInput, ticket)
			if err != nil {
				PrintError("agent run", fmt.Sprintf("Error running agent: %v", err), cmd)
				return err
			}

			if ticket != nil {
				updateTicket(ticket, a, store, verboseFlag)
			}

			if usedTemplate != "" {
				resp = fmt.Sprintf("%s\n%s", usedTemplate, resp)
			}

			out := map[string]interface{}{
				"agentID":    a.ID(),
				"name":       a.Name(),
				"input":      finalInput,
				"response":   resp,
				"parameters": finalParams,
				"model":      modelFlag,
				"ticketID":   ticketFlag,
				"status":     "ok",
			}
			Print(out, "", cmd)
			return nil
		},
	}

	// -------- run flags --------
	runCmd.Flags().StringVarP(&providerFlag, "provider", "p", "venice", "AI provider to use")
	runCmd.Flags().StringVarP(&modelFlag, "model", "m", "default", "Model or endpoint to use")
	runCmd.Flags().StringVar(&cliPromptTemplate, "prompt_template", "", "Override prompt template")
	runCmd.Flags().StringVar(&cliParametersJSON, "parameters", "", "Override parameters JSON")
	runCmd.Flags().StringVar(&ticketFlag, "ticket", "", "Attach an existing workflow ticket ID")
	runCmd.Flags().BoolVar(&verboseFlag, "verbose", false, "Enable verbose ticket step logging")

	agentCmd.AddCommand(runCmd)
	agentCmd.PersistentFlags().Bool("json", false, "Output results in JSON format")

	return agentCmd
}

// ---------------- Helper functions ----------------

func printOrJSON(obj interface{}, msg string, cmd *cobra.Command) {
	jsonFlag := getJSONFlag(cmd)
	if jsonFlag && obj != nil {
		data, _ := json.MarshalIndent(obj, "", "  ")
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
	} else if msg != "" {
		fmt.Fprintln(cmd.OutOrStdout(), msg)
	}
}

// loadOrCreateTicket ensures the ticket exists
func loadOrCreateTicket(ticketID string) (*tickets.Ticket, *tickets.Store, error) {
	if ticketID == "" {
		return nil, nil, nil
	}
	store := tickets.NewStore(tickets.TicketDir)
	tkt, err := store.Load("default", ticketID)
	if err != nil {
		// Create new ticket if missing
		tkt = tickets.NewTicket("default", ticketID, nil)
		if saveErr := store.Save(tkt); saveErr != nil {
			return nil, nil, saveErr
		}
	} else if err := tkt.Validate(); err != nil {
		return nil, nil, err
	}
	return tkt, store, nil
}

func mergeCLIParams(base map[string]string, cliJSON string, cmd *cobra.Command) (map[string]string, error) {
	if cliJSON == "" {
		return base, nil
	}
	var cliParams map[string]string
	if err := json.Unmarshal([]byte(cliJSON), &cliParams); err != nil {
		PrintError("agent run", fmt.Sprintf("Invalid JSON parameters: %v", err), cmd)
		return nil, err
	}
	for k, v := range cliParams {
		base[k] = v
	}
	return base, nil
}

func applyPromptTemplate(agentTemplate, cliTemplate string, params map[string]string, input string) (string, string) {
	usedTemplate := agentTemplate
	if cliTemplate != "" {
		usedTemplate = cliTemplate
	}
	finalInput := input
	if usedTemplate != "" {
		for k, v := range params {
			usedTemplate = strings.ReplaceAll(usedTemplate, fmt.Sprintf("{{%s}}", k), v)
		}
		finalInput = fmt.Sprintf("%s\n%s", usedTemplate, input)
	}
	return finalInput, usedTemplate
}

func updateTicket(ticket *tickets.Ticket, a agent.Agent, store *tickets.Store, verbose bool) {
	ticket.IncrementStep(verbose)
	if ca, ok := a.(agent.ContextualAgent); ok {
		for k, v := range ca.ContextData() {
			ticket.SetNamespaced(a.ID(), k, v)
		}
	}
	store.Save(ticket)
}
