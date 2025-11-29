package tickets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store handles file-based persistence of tickets.
// Each ticket is saved as a JSON file in the Store's directory.
type Store struct {
	dir string     // Directory where tickets are stored
	mu  sync.Mutex // Mutex to protect concurrent access
}

// NewStore creates a new ticket store at the specified directory.
// If dir is empty, it defaults to KEYSTONE_TICKET_DIR or TicketDir.
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

// sanitizeID replaces '/' in ticket IDs to prevent file path issues.
func sanitizeID(id string) string {
	return filepath.Base(strings.ReplaceAll(id, "/", "_"))
}

// ticketPath returns the full file path for a ticket given userID and ticketID.
func (s *Store) ticketPath(userID, ticketID string) string {
	return filepath.Join(s.dir, fmt.Sprintf("%s_%s.json", userID, sanitizeID(ticketID)))
}

// Save writes a Ticket to disk as a JSON file.
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

// Load retrieves a Ticket from disk by userID and ticketID.
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

// List returns all tickets for a specific user, or all users if userID is "all".
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

// Delete removes a specific ticket file for the given userID and ticketID.
func (s *Store) Delete(userID, ticketID string) error {
	return os.Remove(s.ticketPath(userID, ticketID))
}

// Purge deletes all tickets for the specified userID.
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

// Cleanup removes expired or over-hopped tickets for a user and returns the count of removed tickets.
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
