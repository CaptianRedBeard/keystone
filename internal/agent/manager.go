// manager.go
package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"keystone/internal/providers"
	"keystone/internal/providers/venice"

	"gopkg.in/yaml.v3"
)

// AgentManager manages a collection of agents.
type AgentManager struct {
	agents map[string]Agent
	mu     sync.RWMutex
}

// NewManager creates a new empty manager.
func NewManager() *AgentManager {
	return &AgentManager{
		agents: make(map[string]Agent),
	}
}

// Register adds a new agent to the manager.
func (m *AgentManager) Register(a Agent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[a.ID()]; exists {
		return fmt.Errorf("agent with ID %s already exists", a.ID())
	}

	m.agents[a.ID()] = a
	return nil
}

// Get retrieves an agent by its ID.
func (m *AgentManager) Get(id string) (Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, exists := m.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}
	return a, nil
}

// List returns a slice of all registered agents.
func (m *AgentManager) List() []Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]Agent, 0, len(m.agents))
	for _, a := range m.agents {
		list = append(list, a)
	}
	return list
}

// LoadFromConfig loads all YAML agent configs from a directory.
func (m *AgentManager) LoadFromConfig(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to list agent config files: %w", err)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		var cfg struct {
			ID          string `yaml:"id"`
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Provider    string `yaml:"provider"`
			Model       string `yaml:"model"`
			Memory      string `yaml:"memory"`
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("failed to parse YAML in %s: %w", file, err)
		}

		provider, err := createProvider(cfg.Provider)
		if err != nil {
			return fmt.Errorf("failed to init provider for %s: %w", cfg.ID, err)
		}

		agent := NewAgent(cfg.ID, cfg.Name, cfg.Description, provider, cfg.Model, cfg.Memory)
		if err := m.Register(agent); err != nil {
			return fmt.Errorf("failed to register agent %s: %w", cfg.ID, err)
		}
	}

	return nil
}

func createProvider(name string) (providers.Provider, error) {
	switch name {
	case "venice":
		return venice.New("MOCK_API_KEY", ""), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", name)
	}
}

func (m *AgentManager) Unregister(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}
	delete(m.agents, id)
	return nil
}
