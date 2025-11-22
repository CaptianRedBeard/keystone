package cmd

import (
	"context"
	"testing"
	"time"

	"keystone/internal/tickets"
)

// MockAgent implements a minimal agent for testing ticket handling.
type MockAgent struct {
	IDName string
}

func (a *MockAgent) ID() string                    { return a.IDName }
func (a *MockAgent) Name() string                  { return "Mock Agent" }
func (a *MockAgent) Description() string           { return "Mock agent for testing" }
func (a *MockAgent) DefaultModel() string          { return "default" }
func (a *MockAgent) Parameters() map[string]string { return map[string]string{} }
func (a *MockAgent) PromptTemplate() string        { return "" }
func (a *MockAgent) Handle(ctx context.Context, input string, t *tickets.Ticket) (string, error) {
	if t != nil {
		t.SetNamespaced(a.ID(), "last_input", input)
		t.IncrementStep(true)
	}
	return "mock response", nil
}

// -------------------- Tests --------------------

func TestTicketStepIncrement(t *testing.T) {
	agent := &MockAgent{IDName: "mock"}
	ticket := tickets.NewTicket("t1", "user1", nil)

	if ticket.Step != 0 || ticket.Hops != 0 {
		t.Fatalf("expected initial step/hops 0, got %d/%d", ticket.Step, ticket.Hops)
	}

	_, err := agent.Handle(context.Background(), "input", ticket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ticket.Step != 1 || ticket.Hops != 1 {
		t.Errorf("expected step/hops 1/1, got %d/%d", ticket.Step, ticket.Hops)
	}

	ttl := ticket.ExpiresAt.Sub(ticket.CreatedAt).Round(time.Second)
	if ttl != tickets.DefaultTTL {
		t.Errorf("expected TTL %v, got %v", tickets.DefaultTTL, ttl)
	}
}

func TestTicketExpired(t *testing.T) {
	ticket := tickets.NewTicket("t2", "user2", nil)
	ticket.ExpiresAt = time.Now().Add(-time.Minute)

	if err := ticket.Validate(); err == nil {
		t.Fatal("expected error for expired ticket, got none")
	}
}

func TestTicketMaxHops(t *testing.T) {
	ticket := tickets.NewTicket("t3", "user3", nil)
	ticket.MaxHops = 2
	ticket.Hops = 2

	if err := ticket.Validate(); err == nil {
		t.Fatal("expected error for ticket reaching max hops, got none")
	}
}

func TestAgentRunWithoutTicket(t *testing.T) {
	agent := &MockAgent{IDName: "mock"}
	_, err := agent.Handle(context.Background(), "input", nil)
	if err != nil {
		t.Fatalf("unexpected error running agent without ticket: %v", err)
	}
}

func TestTicketLifecycleWithContext(t *testing.T) {
	agentA := &MockAgent{IDName: "agentA"}
	agentB := &MockAgent{IDName: "agentB"}
	ticket := tickets.NewTicket("t4", "user4", map[string]string{"agentA.foo": "bar"})

	_, _ = agentA.Handle(context.Background(), "helloA", ticket)
	_, _ = agentB.Handle(context.Background(), "helloB", ticket)

	aCtx := ticket.GetAllNamespaced("agentA")
	if val := aCtx["last_input"]; val == "helloB" {
		t.Errorf("agentA should not see agentB key")
	}
	if val := aCtx["foo"]; val != "bar" {
		t.Errorf("expected agentA.foo='bar', got '%s'", val)
	}

	bCtx := ticket.GetAllNamespaced("agentB")
	if _, exists := bCtx["foo"]; exists {
		t.Errorf("agentB should not see agentA original key 'foo'")
	}
}

func TestAgentCannotOverwriteOtherNamespace(t *testing.T) {
	agentA := &MockAgent{IDName: "agentA"}
	agentB := &MockAgent{IDName: "agentB"}
	ticket := tickets.NewTicket("t5", "user5", map[string]string{"agentA.secret": "123"})

	_, _ = agentA.Handle(context.Background(), "helloA", ticket)
	_, _ = agentB.Handle(context.Background(), "new_value", ticket)

	if val := ticket.GetAllNamespaced("agentA")["secret"]; val != "123" {
		t.Errorf("agentA's secret was overwritten by agentB, got '%s'", val)
	}
	if val := ticket.GetAllNamespaced("agentB")["last_input"]; val != "new_value" {
		t.Errorf("agentB did not correctly write to its namespace, got '%s'", val)
	}
}
