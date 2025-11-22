package usage

import (
	"sync"
	"testing"
	"time"
)

func TestUsageTracker(t *testing.T) {
	tracker := NewTracker()

	t.Run("RecordEntry", func(t *testing.T) {
		e := tracker.Record("agent1", "venice", 10)
		if e.AgentID != "agent1" {
			t.Errorf("expected AgentID 'agent1', got '%s'", e.AgentID)
		}
		if e.Tokens != 10 {
			t.Errorf("expected Tokens 10, got %d", e.Tokens)
		}
		if time.Since(e.Timestamp) > time.Second {
			t.Error("Timestamp seems incorrect")
		}
	})

	t.Run("SummaryAndList", func(t *testing.T) {
		tracker.Record("agent1", "venice", 5)
		summary := tracker.Summary()
		if summary.TotalRequests != 2 {
			t.Errorf("expected total requests 2, got %d", summary.TotalRequests)
		}
		if summary.TotalTokens != 15 {
			t.Errorf("expected total tokens 15, got %d", summary.TotalTokens)
		}

		list := tracker.List()
		if len(list) != 2 {
			t.Errorf("expected 2 entries in list, got %d", len(list))
		}
	})
}

func TestGenerateIDUniqueness(t *testing.T) {
	id1 := generateID()
	time.Sleep(time.Microsecond) // ensure a tiny time gap
	id2 := generateID()

	if id1 == id2 {
		t.Errorf("expected unique IDs, got same: %s", id1)
	}
}

func TestConcurrentRecordSafety(t *testing.T) {
	tracker := NewTracker()
	wg := sync.WaitGroup{}
	agents := []string{"a", "b", "c", "d"}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tracker.Record(agents[i%len(agents)], "test_provider", i)
		}(i)
	}
	wg.Wait()

	summary := tracker.Summary()
	if summary.TotalRequests != 100 {
		t.Errorf("expected 100 total requests, got %d", summary.TotalRequests)
	}
	if summary.TotalTokens == 0 {
		t.Error("expected non-zero total tokens after concurrent writes")
	}
}
