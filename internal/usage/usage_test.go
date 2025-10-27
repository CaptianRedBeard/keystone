package usage

import (
	"testing"
	"time"
)

func TestUsageTracker(t *testing.T) {
	tracker := NewTracker()

	// Record first entry
	e1 := tracker.Record("agent1", "venice", 10)
	if e1.AgentID != "agent1" {
		t.Errorf("expected AgentID 'agent1', got '%s'", e1.AgentID)
	}
	if e1.Tokens != 10 {
		t.Errorf("expected Tokens 10, got %d", e1.Tokens)
	}
	if time.Since(e1.Timestamp) > time.Second {
		t.Error("Timestamp seems incorrect")
	}

	// Record another entry
	tracker.Record("agent1", "venice", 5)
	summary := tracker.Summary()
	if summary.TotalTokens != 15 {
		t.Errorf("expected total tokens 15, got %d", summary.TotalTokens)
	}
	if summary.TotalRequests != 2 {
		t.Errorf("expected total requests 2, got %d", summary.TotalRequests)
	}

	// Test List
	list := tracker.List()
	if len(list) != 2 {
		t.Errorf("expected 2 entries in list, got %d", len(list))
	}
}
