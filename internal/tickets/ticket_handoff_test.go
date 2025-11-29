package tickets

import (
	"testing"
	"time"
)

func TestTicketHandoffBasic(t *testing.T) {
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

	// Check namespaced context initialized for agent
	key := Namespaced("agent1", "_")
	if _, ok := ticket.Context[key]; !ok {
		t.Errorf("expected %q context to be initialized", key)
	}
}

func TestTicketHandoffMaxHops(t *testing.T) {
	ticket := NewTicket("t2", "user1", nil)
	ticket.Hops = ticket.MaxHops

	err := ticket.Handoff("agent1")
	if err == nil {
		t.Fatal("expected error for exceeding max hops")
	}
}

func TestTicketHandoffExpired(t *testing.T) {
	ticket := NewTicket("t3", "user1", nil)
	ticket.ExpiresAt = time.Now().Add(-time.Minute)

	err := ticket.Handoff("agent1")
	if err == nil {
		t.Fatal("expected error for expired ticket")
	}
}

func TestTicketHandoffHook(t *testing.T) {
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

	// Check namespaced context initialized for agent
	key := Namespaced("agentX", "_")
	if _, ok := ticket.Context[key]; !ok {
		t.Errorf("expected %q context to be initialized", key)
	}
}
