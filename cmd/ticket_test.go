package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"keystone/internal/agent"
	"keystone/internal/config"
	"keystone/internal/tickets"
)

func TestTicketCLI(t *testing.T) {
	tmpDir := t.TempDir()
	store := tickets.NewStore(tmpDir)

	// Dummy agent provider
	dummyProvider := func(dir string) *agent.AgentManager { return agent.NewManager() }

	// Helper: create a ticket in the store
	createTicket := func() *tickets.Ticket {
		tkt := tickets.NewTicket(tickets.NewID("cli", "ticket", "default"), "default", nil)
		tkt.MaxHops = 5
		tkt.TTL = time.Hour
		if err := store.Save(tkt); err != nil {
			t.Fatalf("failed to save ticket: %v", err)
		}
		return tkt
	}

	// Helper: run a CLI command and capture output
	runCommand := func(args ...string) string {
		buf := new(bytes.Buffer)
		cfgLoader := func(_ string) (*config.Config, error) { return config.New(), nil }
		cmd := NewRootCmd(dummyProvider, cfgLoader, buf, store)
		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("command %v failed: %v", args, err)
		}
		return buf.String()
	}

	// Helper: parse JSON output
	parseJSONOutput := func(output string, t *testing.T) map[string]interface{} {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
		}
		return data
	}

	// ------------------------ Tests ------------------------

	// ticket new --json
	output := runCommand("ticket", "new", "--json")
	ticketData := parseJSONOutput(output, t)
	ticketID, ok := ticketData["id"].(string)
	if !ok || ticketID == "" {
		t.Fatal("failed to extract ticket ID from 'ticket new'")
	}

	// ticket list --json
	output = runCommand("ticket", "list", "--json")
	var listData []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &listData); err != nil {
		t.Fatalf("failed to unmarshal JSON list output: %v", err)
	}
	if len(listData) == 0 {
		t.Errorf("expected at least 1 ticket in list, got 0")
	}

	// ticket cleanup --json with expired ticket
	expiredTicket := createTicket()
	expiredTicket.ExpiresAt = expiredTicket.CreatedAt.Add(-time.Hour)
	if err := store.Save(expiredTicket); err != nil {
		t.Fatalf("failed to save expired ticket: %v", err)
	}

	output = runCommand("ticket", "cleanup", "--json")
	cleanupData := parseJSONOutput(output, t)
	if removed, ok := cleanupData["removed"].(float64); !ok || removed < 1 {
		t.Errorf("expected at least 1 removed ticket, got %v", cleanupData["removed"])
	}

	// ticket monitor --json
	createTicket() // fresh ticket
	staleTicket := createTicket()
	staleTicket.Hops = staleTicket.MaxHops + 1
	if err := store.Save(staleTicket); err != nil {
		t.Fatalf("failed to save stale ticket: %v", err)
	}

	output = runCommand("ticket", "monitor", "--json")
	monitorData := parseJSONOutput(output, t)
	if monitorData["total_tickets"].(float64) < 2 {
		t.Errorf("expected at least 2 tickets in monitor, got %v", monitorData["total_tickets"])
	}
	if monitorData["stale_tickets"].(float64) < 1 {
		t.Errorf("expected at least 1 stale ticket, got %v", monitorData["stale_tickets"])
	}
}
