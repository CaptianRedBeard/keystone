// Package agent provides base implementations and helpers for AI agents.
package agent

import (
	"fmt"
	"sync"
)

// AgentManager manages a thread-safe collection of Agents.
type AgentManager struct {
	agents map[string]Agent
	mu     sync.RWMutex
}

// NewManager creates and returns a new AgentManager.
func NewManager() *AgentManager {
	return &AgentManager{agents: make(map[string]Agent)}
}

// Register adds a new agent to the manager.
func (m *AgentManager) Register(a Agent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.agents[a.ID()]; ok {
		return fmt.Errorf("agent %s already exists", a.ID())
	}
	m.agents[a.ID()] = a
	return nil
}

// Get retrieves an agent by ID.
func (m *AgentManager) Get(id string) (Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.agents[id]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", id)
	}
	return a, nil
}

// List returns a slice of all registered agents.
func (m *AgentManager) List() []Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Agent, 0, len(m.agents))
	for _, a := range m.agents {
		out = append(out, a)
	}
	return out
}

// Unregister removes an agent by ID.
func (m *AgentManager) Unregister(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.agents[id]; !ok {
		return fmt.Errorf("agent %s not found", id)
	}
	delete(m.agents, id)
	return nil
}
