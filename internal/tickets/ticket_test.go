package tickets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// setupStore creates a temporary ticket store for tests and returns a cleanup func.
func setupStore(t *testing.T) (*Store, func()) {
	t.Helper()
	tmpDir := filepath.Join(os.TempDir(), "keystone_tickets_test")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatal(err)
	}
	store := NewStore(tmpDir)
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return store, cleanup
}

// TestTicketLifecycle contains core ticket persistence and utility tests.
func TestTicketLifecycle(t *testing.T) {
	store, cleanup := setupStore(t)
	defer cleanup()
	userID := "user1"

	t.Run("Create_and_validate_ticket", func(t *testing.T) {
		ticket := NewTicket("ticket1", userID, map[string]string{"foo": "bar"})
		if ticket.ID == "" {
			t.Error("expected ticket ID to be set")
		}
		if err := ticket.Validate(); err != nil {
			t.Errorf("expected ticket to be valid, got %v", err)
		}
	})

	t.Run("Save_and_load_ticket", func(t *testing.T) {
		ticket := NewTicket("ticket2", userID, nil)
		if err := store.Save(ticket); err != nil {
			t.Fatalf("failed to save ticket: %v", err)
		}
		loaded, err := store.Load(userID, ticket.ID)
		if err != nil {
			t.Fatalf("failed to load ticket: %v", err)
		}
		if loaded.ID != ticket.ID {
			t.Errorf("expected ID %s, got %s", ticket.ID, loaded.ID)
		}
	})

	t.Run("List_tickets", func(t *testing.T) {
		tickets, err := store.List(userID)
		if err != nil {
			t.Fatal(err)
		}
		if len(tickets) == 0 {
			t.Error("expected at least 1 ticket")
		}
	})

	t.Run("Delete_ticket", func(t *testing.T) {
		ticket := NewTicket("ticket3", userID, nil)
		_ = store.Save(ticket)
		if err := store.Delete(userID, ticket.ID); err != nil {
			t.Fatal(err)
		}
		list, _ := store.List(userID)
		for _, tck := range list {
			if tck.ID == ticket.ID {
				t.Error("ticket not deleted")
			}
		}
	})

	t.Run("Purge_tickets", func(t *testing.T) {
		t1 := NewTicket("ticket4", userID, nil)
		t2 := NewTicket("ticket5", "other", nil)
		_ = store.Save(t1)
		_ = store.Save(t2)
		if err := store.Purge(userID); err != nil {
			t.Fatal(err)
		}
		list, _ := store.List(userID)
		if len(list) != 0 {
			t.Error("expected 0 tickets for user after purge")
		}
	})

	t.Run("Cleanup_expired_or_over_hopped_tickets", func(t *testing.T) {
		ticket := NewTicket("ticket6", userID, nil)
		ticket.ExpiresAt = time.Now().Add(-time.Minute)
		ticket.Hops = ticket.MaxHops
		_ = store.Save(ticket)
		removed, err := store.Cleanup(userID)
		if err != nil {
			t.Fatal(err)
		}
		if removed != 1 {
			t.Errorf("expected 1 removed ticket, got %d", removed)
		}
	})

	t.Run("Namespaced_context", func(t *testing.T) {
		ticket := NewTicket("ticket7", userID, nil)
		ticket.SetNamespaced("agent1", "key", "value")
		v, ok := ticket.GetNamespaced("agent1", "key")
		if !ok || v != "value" {
			t.Error("expected namespaced value 'value'")
		}
		all := ticket.GetAllNamespaced("agent1")
		if all["key"] != "value" {
			t.Error("expected GetAllNamespaced to include 'key'")
		}
	})

	t.Run("IncrementStep_updates_Step_and_Hops", func(t *testing.T) {
		ticket := NewTicket("ticket8", userID, nil)
		ticket.IncrementStep(false)
		if ticket.Step != 1 || ticket.Hops != 1 {
			t.Errorf("expected Step=1, Hops=1, got Step=%d, Hops=%d", ticket.Step, ticket.Hops)
		}
	})

	t.Run("Serialize_and_NewID", func(t *testing.T) {
		ticket := NewTicket("ticket9", userID, nil)
		data := ticket.Serialize()
		if data["id"] != ticket.ID {
			t.Errorf("expected Serialize id %s, got %s", ticket.ID, data["id"])
		}

		id := NewID("part1", "part2")
		if !strings.Contains(id, "part1") || !strings.Contains(id, "part2") {
			t.Errorf("expected NewID to contain parts, got %s", id)
		}
	})
}

// TestTicketHandoff contains Phase 4 tests for handoff behavior (step/hops, TTL, hooks).
func TestTicketHandoff(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		ticket := NewTicket("t1", "user1", nil)
		initialStep := ticket.Step
		initialHops := ticket.Hops

		err := ticket.Handoff("agent1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if ticket.Step != initialStep+1 {
			t.Errorf("expected Step %d, got %d", initialStep+1, ticket.Step)
		}
		if ticket.Hops != initialHops+1 {
			t.Errorf("expected Hops %d, got %d", initialHops+1, ticket.Hops)
		}

		// Check context initialized for agent
		if _, ok := ticket.Context["agent1"]; !ok {
			t.Error("expected agent1 context to be initialized")
		}
	})

	t.Run("MaxHops", func(t *testing.T) {
		ticket := NewTicket("t2", "user1", nil)
		ticket.Hops = ticket.MaxHops

		err := ticket.Handoff("agent1")
		if err == nil {
			t.Fatal("expected error for exceeding max hops")
		}
	})

	t.Run("Expired", func(t *testing.T) {
		ticket := NewTicket("t3", "user1", nil)
		ticket.ExpiresAt = time.Now().Add(-time.Minute)

		err := ticket.Handoff("agent1")
		if err == nil {
			t.Fatal("expected error for expired ticket")
		}
	})

	t.Run("Hook", func(t *testing.T) {
		ticket := NewTicket("t4", "user1", nil)

		called := false
		OnHandoffHook = func(tk *Ticket, agentID string) {
			called = true
			if tk.ID != ticket.ID {
				t.Errorf("expected ticket ID %s, got %s", ticket.ID, tk.ID)
			}
			if agentID != "agentX" {
				t.Errorf("expected agentID 'agentX', got %s", agentID)
			}
		}
		defer func() { OnHandoffHook = nil }() // reset after test

		err := ticket.Handoff("agentX")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !called {
			t.Error("expected OnHandoffHook to be called")
		}
	})
}
