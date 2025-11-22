package tickets

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"keystone/internal/logger"
)

// Default ticket storage path, can be overridden with KEYSTONE_TICKET_DIR
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

// OnHandoffHook is an optional function to receive handoff events (Phase 5+ GUI/monitoring).
var OnHandoffHook func(t *Ticket, nextAgentID string)

// NewTicket constructs a new Ticket with optional context.
func NewTicket(id, userID string, ctx interface{}) *Ticket {
	var contextMap map[string]interface{}
	switch v := ctx.(type) {
	case map[string]string:
		contextMap = make(map[string]interface{}, len(v))
		for k, val := range v {
			contextMap[k] = val
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

// NewID generates a unique ticket ID using optional string parts and current timestamp.
func NewID(parts ...string) string {
	return fmt.Sprintf("%s-%d", filepath.Join(parts...), time.Now().UnixNano())
}

// Validate checks if the ticket is expired or exceeded max hops.
func (t *Ticket) Validate() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if time.Now().After(t.ExpiresAt) {
		return errors.New("ticket expired")
	}
	if t.Hops >= t.MaxHops {
		return fmt.Errorf("ticket exceeded max hops (%d)", t.MaxHops)
	}
	return nil
}

// IsExpired returns true if the ticket has expired.
func (t *Ticket) IsExpired() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return time.Now().After(t.ExpiresAt)
}

// GetNamespaced retrieves a namespaced value for a specific agent.
func (t *Ticket) GetNamespaced(agentID, key string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if val, ok := t.Context[agentID+"."+key]; ok {
		if s, ok := val.(string); ok {
			return s, true
		}
	}
	return "", false
}

// SetNamespaced sets a namespaced value for a specific agent.
func (t *Ticket) SetNamespaced(agentID, key, value string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Context[agentID+"."+key] = value
}

// GetAllNamespaced returns all key-value pairs for a specific agent.
func (t *Ticket) GetAllNamespaced(agentID string) map[string]string {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[string]string)
	prefix := agentID + "."
	for k, v := range t.Context {
		if s, ok := v.(string); ok && strings.HasPrefix(k, prefix) {
			out[k[len(prefix):]] = s
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
// It validates TTL/MaxHops, increments step and hops, and ensures
// the next agent has an initialized context namespace.
// Returns an error if the ticket cannot be handed off.
func (t *Ticket) Handoff(nextAgentID string) error {
	// Validate using existing method (which handles locking).
	if err := t.Validate(); err != nil {
		return fmt.Errorf("cannot handoff ticket %s: %w", t.ID, err)
	}

	// Increment step/hops using the existing method (which handles locking and logging).
	t.IncrementStep(true)

	// Initialize the next agent's namespace in the context in a single locked section.
	t.mu.Lock()
	if _, exists := t.Context[nextAgentID]; !exists {
		t.Context[nextAgentID] = make(map[string]interface{})
	}
	// capture current step/hops for logging without holding lock longer than needed
	step := t.Step
	hops := t.Hops
	t.mu.Unlock()

	// Log the handoff (do not hold mutex while logging or calling hooks).
	logger.Info(fmt.Sprintf("Ticket %s handed off to agent %s (Step=%d, Hops=%d)",
		t.ID, nextAgentID, step, hops), false)

	// Call hook outside of locks so hooks can't deadlock by calling ticket methods.
	if OnHandoffHook != nil {
		OnHandoffHook(t, nextAgentID)
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

// -------------------- Store --------------------

// Store handles file-based persistence of tickets.
type Store struct {
	dir string
	mu  sync.Mutex
}

// NewStore creates a new ticket store at the given directory.
func NewStore(dir string) *Store {
	if dir == "" {
		if envDir := os.Getenv("KEYSTONE_TICKET_DIR"); envDir != "" {
			dir = envDir
		} else {
			dir = TicketDir
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(fmt.Sprintf("failed to create ticket directory %s: %v", dir, err))
	}
	return &Store{dir: dir}
}

func sanitizeID(id string) string {
	return filepath.Base(strings.ReplaceAll(id, "/", "_"))
}

func (s *Store) ticketPath(userID, ticketID string) string {
	return filepath.Join(s.dir, fmt.Sprintf("%s_%s.json", userID, sanitizeID(ticketID)))
}

// Save persists a ticket to disk.
func (s *Store) Save(t *Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fullPath := s.ticketPath(t.UserID, t.ID)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fullPath, data, 0o644)
}

// Load retrieves a ticket by user ID and ticket ID.
func (s *Store) Load(userID, ticketID string) (*Ticket, error) {
	data, err := os.ReadFile(s.ticketPath(userID, ticketID))
	if err != nil {
		return nil, err
	}
	var t Ticket
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// List returns all tickets for a user or all users if "all" is passed.
func (s *Store) List(userID string) ([]*Ticket, error) {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var tickets []*Ticket
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.dir, f.Name()))
		if err != nil {
			continue
		}
		var t Ticket
		if err := json.Unmarshal(data, &t); err != nil {
			continue
		}
		if userID == "all" || t.UserID == userID {
			tickets = append(tickets, &t)
		}
	}
	return tickets, nil
}

// Delete removes a ticket from storage.
func (s *Store) Delete(userID, ticketID string) error {
	return os.Remove(s.ticketPath(userID, ticketID))
}

// Purge deletes all tickets for a given user.
func (s *Store) Purge(userID string) error {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(s.dir, f.Name()))
		if err != nil {
			continue
		}
		var t Ticket
		if err := json.Unmarshal(data, &t); err != nil {
			continue
		}
		if t.UserID == userID {
			_ = os.Remove(filepath.Join(s.dir, f.Name()))
		}
	}
	return nil
}

// Cleanup removes expired or over-hopped tickets for a user, returning count of removed tickets.
func (s *Store) Cleanup(userID string) (int, error) {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return 0, err
	}

	removed := 0
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(s.dir, f.Name()))
		if err != nil {
			continue
		}
		var t Ticket
		if err := json.Unmarshal(data, &t); err != nil {
			continue
		}
		if t.UserID != userID {
			continue
		}
		if t.IsExpired() || t.Hops >= t.MaxHops {
			_ = os.Remove(filepath.Join(s.dir, f.Name()))
			removed++
		}
	}
	return removed, nil
}
