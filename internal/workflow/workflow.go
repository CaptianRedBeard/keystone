package workflow

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// -------------------------
// Defaults
// -------------------------

var DefaultWorkflowDir = filepath.Join(".", "workflows")

// -------------------------
// Workflow Store
// -------------------------

type Store struct {
	dir string
	mu  sync.Mutex
}

func NewStore(dir string) *Store {
	if dir == "" {
		dir = DefaultWorkflowDir
	}
	_ = os.MkdirAll(dir, 0o755)
	return &Store{dir: dir}
}

func sanitizeID(id string) string {
	return filepath.Base(strings.ReplaceAll(id, "/", "_"))
}

func (s *Store) workflowPath(id string) string {
	return filepath.Join(s.dir, sanitizeID(id)+".json")
}

// Save persists a workflow
func (s *Store) Save(wf Workflow) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.workflowPath(wf.ID), data, 0o644)
}

// Load retrieves a workflow
func (s *Store) Load(id string) (*Workflow, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.workflowPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}

	var wf Workflow
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, err
	}
	return &wf, nil
}

// List all workflows
func (s *Store) List() ([]Workflow, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	files, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	workflows := []Workflow{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(s.dir, f.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var wf Workflow
		if err := json.Unmarshal(data, &wf); err != nil {
			continue
		}
		workflows = append(workflows, wf)
	}
	return workflows, nil
}

// Delete removes a workflow by ID
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return os.Remove(s.workflowPath(id))
}

// -------------------------
// Workflow structs (existing)
// -------------------------

type Step struct {
	AgentID string            `yaml:"agent_id"`
	Input   string            `yaml:"input"`
	Params  map[string]string `yaml:"params,omitempty"`
}

type Workflow struct {
	ID    string `yaml:"id"`
	Steps []Step `yaml:"steps"`
}

type StepResult struct {
	AgentID string
	Output  string
	Error   error
}
