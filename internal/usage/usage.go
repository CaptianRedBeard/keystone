package usage

import (
	"sync"
	"time"
)

// Entry represents a single usage event.
type Entry struct {
	RequestID string    // unique ID per request
	AgentID   string    // which agent made the request
	Provider  string    // provider used (e.g., venice)
	Tokens    int       // tokens used
	Timestamp time.Time // time of the request
}

// Tracker keeps track of usage events.
type Tracker struct {
	mu      sync.Mutex
	entries []Entry
}

// Summary holds aggregated usage info.
type Summary struct {
	TotalRequests int
	TotalTokens   int
}

// NewTracker creates a usage tracker instance.
func NewTracker() *Tracker {
	return &Tracker{
		entries: make([]Entry, 0),
	}
}

// Record adds a new usage entry to the tracker.
func (t *Tracker) Record(agentID, provider string, tokens int) Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	e := Entry{
		RequestID: generateID(),
		AgentID:   agentID,
		Provider:  provider,
		Tokens:    tokens,
		Timestamp: time.Now(),
	}

	t.entries = append(t.entries, e)
	return e
}

// Summary aggregates usage data.
func (t *Tracker) Summary() Summary {
	t.mu.Lock()
	defer t.mu.Unlock()

	s := Summary{}
	for _, e := range t.entries {
		s.TotalRequests++
		s.TotalTokens += e.Tokens
	}
	return s
}

// List returns a copy of all usage entries.
func (t *Tracker) List() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	copied := make([]Entry, len(t.entries))
	copy(copied, t.entries)
	return copied
}

// generateID produces a simple timestamp-based unique ID.
func generateID() string {
	return time.Now().Format("20060102-150405.000000")
}
