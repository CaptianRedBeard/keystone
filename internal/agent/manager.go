package agent

import (
	"fmt"
	"sync"
)

type AgentManager struct {
	agents map[string]*Agent
	mu     sync.RWMutex
}

func NewManager() *AgentManager {
	return &AgentManager{
		agents: make(map[string]*Agent),
	}
}

func (m *AgentManager) Register(a *Agent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[a.ID]; exists {
		return fmt.Errorf("agent with ID %s already exists", a.ID)
	}

	m.agents[a.ID] = a
	return nil
}

func (m *AgentManager) Get(id string) (*Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, exists := m.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}
	return a, nil
}

func (m *AgentManager) List() []*Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*Agent, 0, len(m.agents))
	for _, a := range m.agents {
		list = append(list, a)
	}
	return list
}
