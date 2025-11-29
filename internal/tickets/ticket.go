package tickets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"keystone/internal/logger"
)

// Default ticket storage path
var TicketDir = filepath.Join(os.Getenv("HOME"), ".keystone", "tickets")

const (
	DefaultTTL     = time.Hour
	DefaultMaxHops = 5
)

// Ticket represents a unit of work tracked across agents.
type Ticket struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Context   map[string]interface{} `json:"context"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
	Hops      int                    `json:"hops"`
	Step      int                    `json:"step"`
	MaxHops   int                    `json:"max_hops"`
	TTL       time.Duration          `json:"ttl"`
	mu        sync.Mutex
}

// OnHandoffHook is an optional function to receive handoff events.
var OnHandoffHook func(t *Ticket, nextAgentID string)

// NewTicket constructs a new Ticket with optional context.
func NewTicket(id, userID string, ctx interface{}) *Ticket {
	var contextMap map[string]interface{}
	switch v := ctx.(type) {
	case map[string]string:
		contextMap = make(map[string]interface{}, len(v))
		for k, val := range v {
			contextMap[Namespaced("", k)] = val
		}
	case map[string]interface{}:
		contextMap = v
	case nil:
		contextMap = make(map[string]interface{})
	default:
		panic(fmt.Sprintf("unsupported context type %T", ctx))
	}

	now := time.Now()
	return &Ticket{
		ID:        id,
		UserID:    userID,
		Context:   contextMap,
		CreatedAt: now,
		ExpiresAt: now.Add(DefaultTTL),
		MaxHops:   DefaultMaxHops,
		TTL:       DefaultTTL,
	}
}

// GetNamespaced retrieves a namespaced value for a specific agent.
func (t *Ticket) GetNamespaced(agentID, key string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	ns := Namespaced(agentID, key)
	if val, ok := t.Context[ns]; ok {
		if s, ok := val.(string); ok {
			return s, true
		}
	}
	return "", false
}

// SetNamespaced sets a namespaced value for a specific agent (allows overwrite)
func (t *Ticket) SetNamespaced(agentID, key, value string) {
	_ = t.SetNamespacedWithOverwrite(agentID, key, value, true)
}

// SetNamespacedWithOverwrite sets a namespaced value for a specific agent.
func (t *Ticket) SetNamespacedWithOverwrite(agentID, key, value string, allowOverwrite bool) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	ns := Namespaced(agentID, key)

	if !allowOverwrite {
		if _, exists := t.Context[ns]; exists {
			return fmt.Errorf("namespaced key %q already exists for agent %q", key, agentID)
		}
	}

	t.Context[ns] = value
	return nil
}

// GetAllNamespaced returns all key-value pairs for a specific agent.
func (t *Ticket) GetAllNamespaced(agentID string) map[string]string {
	t.mu.Lock()
	defer t.mu.Unlock()

	out := make(map[string]string)
	prefix := "agent." + agentID + "."
	for k, v := range t.Context {
		if s, ok := v.(string); ok && strings.HasPrefix(k, prefix) {
			keyWithoutPrefix := strings.TrimPrefix(k, prefix)
			out[keyWithoutPrefix] = s
		}
	}
	return out
}

// IncrementStep increments the ticket's step and hops, optionally logging.
func (t *Ticket) IncrementStep(verbose bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Step++
	t.Hops++
	if verbose {
		logger.Info(fmt.Sprintf("Ticket %s incremented: Step=%d, Hops=%d", t.ID, t.Step, t.Hops), false)
	}
}

// Handoff transfers the ticket to the next agent.
func (t *Ticket) Handoff(agentID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("ticket expired")
	}
	if t.Hops >= t.MaxHops {
		return fmt.Errorf("ticket max hops exceeded")
	}

	t.Step++
	t.Hops++

	if t.Context == nil {
		t.Context = make(map[string]interface{})
	}
	key := Namespaced(agentID, "_")
	if _, ok := t.Context[key]; !ok {
		t.Context[key] = ""
	}

	if OnHandoffHook != nil {
		OnHandoffHook(t, agentID)
	}

	return nil
}

// Serialize returns a snapshot of ticket fields.
func (t *Ticket) Serialize() map[string]interface{} {
	t.mu.Lock()
	defer t.mu.Unlock()
	return map[string]interface{}{
		"id":         t.ID,
		"user_id":    t.UserID,
		"created_at": t.CreatedAt,
		"expires_at": t.ExpiresAt,
		"hops":       t.Hops,
		"step":       t.Step,
		"max_hops":   t.MaxHops,
		"ttl":        t.TTL.Seconds(),
	}
}

// SerializeContext returns a copy of the ticket's context.
func (t *Ticket) SerializeContext() map[string]interface{} {
	t.mu.Lock()
	defer t.mu.Unlock()
	copyCtx := make(map[string]interface{}, len(t.Context))
	for k, v := range t.Context {
		copyCtx[k] = v
	}
	return copyCtx
}

// IsExpired returns true if the ticket's ExpiresAt time is in the past.
func (t *Ticket) IsExpired() bool {
	t.mu.Lock()
	expires := t.ExpiresAt
	t.mu.Unlock()
	return time.Now().After(expires)
}

// Validate checks if the ticket is expired or has exceeded max hops.
func (t *Ticket) Validate() error {
	t.mu.Lock()
	expires := t.ExpiresAt
	hops := t.Hops
	maxHops := t.MaxHops
	t.mu.Unlock()

	if time.Now().After(expires) {
		return fmt.Errorf("ticket expired")
	}
	if hops >= maxHops {
		return fmt.Errorf("ticket exceeded max hops (%d)", maxHops)
	}
	return nil
}

// NewID generates a unique ticket ID using optional string parts and the current timestamp.
func NewID(parts ...string) string {
	return fmt.Sprintf("%s-%d", filepath.Join(parts...), time.Now().UnixNano())
}
