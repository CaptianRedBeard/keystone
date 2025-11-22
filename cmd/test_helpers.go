package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"keystone/internal/agent"
	"keystone/internal/tickets"

	"github.com/spf13/cobra"
)

// managerProvider returns a new AgentManager for tests
func managerProvider() *agent.AgentManager {
	return agent.NewManager()
}

// captureOutput runs a cobra command and captures its stdout/stderr.
func captureOutput(f func(cmd *cobra.Command)) string {
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	f(cmd)
	return buf.String()
}

// parseJSONOutput trims and parses CLI JSON output
func parseJSONOutput(output string, t *testing.T) map[string]interface{} {
	output = strings.TrimSpace(output)
	if output == "" {
		t.Fatalf("expected JSON output but got empty string")
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(output), &m); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v\nOutput:\n%s", err, output)
	}
	return m
}

// runCLITicketStep simulates a CLI agent run that increments a ticket.
func runCLITicketStep(tmpDir, ticketID string) error {
	store := tickets.NewStore(tmpDir)

	tkt, err := store.Load("default", ticketID)
	if err != nil {
		// Ticket doesn't exist → create it with an empty map for Data
		tkt = tickets.NewTicket("default", ticketID, map[string]interface{}{})
		if err := store.Save(tkt); err != nil {
			return err
		}
	}

	// Increment step/hops
	tkt.Step++
	tkt.Hops++

	return store.Save(tkt)
}
