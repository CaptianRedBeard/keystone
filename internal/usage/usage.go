package usage

import (
	"sync"
	"time"
)

type Entry struct {
	RequestID string    // optional unique ID per request
	AgentID   string    // which agent made the request
	Provider  string    // provider used (e.g., venice)
	Tokens    int       // tokens used
	Timestamp time.Time // time of the request
}

// Tracker keeps track of all usage events.
type Tracker struct {
	mu      sync.Mutex
	entries []Entry
}

type Summary struct {
	TotalRequests int
	TotalTokens   int
}

func NewTracker() *Tracker {
	return &Tracker{
		entries: make([]Entry, 0),
	}
}

// Record adds a new usage entry.
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

// Summary returns total requests and total tokens.
func (t *Tracker) Summary() Summary {
	t.mu.Lock()
	defer t.mu.Unlock()

	total := Summary{}
	for _, e := range t.entries {
		total.TotalRequests++
		total.TotalTokens += e.Tokens
	}
	return total
}

// List returns a copy of all entries
func (t *Tracker) List() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	copied := make([]Entry, len(t.entries))
	copy(copied, t.entries)
	return copied
}

// simple unique ID generator (Phase 1)
func generateID() string {
	return time.Now().Format("20060102-150405.000000")
}
