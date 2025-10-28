// internal/agent/lifecycle.go
package agent

import (
	"fmt"
	"keystone/internal/providers"
	"keystone/internal/providers/venice"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// LifecycleManager extends AgentManager with load/unload/reload support.
type LifecycleManager struct {
	manager   *AgentManager
	mu        sync.Mutex
	configDir string
}

// NewLifecycleManager creates a new manager for a given config directory.
func NewLifecycleManager(configDir string) *LifecycleManager {
	if configDir == "" {
		configDir = "./internal/agent/config"
	}
	return &LifecycleManager{
		manager:   NewManager(),
		configDir: configDir,
	}
}

// BuildAgent constructs an AgentBase from an AgentConfig.
func BuildAgent(cfg AgentConfig) *AgentBase {
	// Select provider based on cfg.Provider
	var provider providers.Provider
	switch cfg.Provider {
	case "venice":
		provider = venice.New("MOCK_API_KEY", "")
	default:
		// fallback or panic; could log warning
		provider = venice.New("MOCK_API_KEY", "")
	}

	return NewAgent(
		cfg.ID,
		cfg.Name,
		cfg.Description,
		provider,
		cfg.Model,
		cfg.Memory,
		WithPromptTemplate(cfg.PromptTemplate),
		WithParameters(cfg.Parameters),
		WithLogging(cfg.Logging),
	)
}

// LoadAgent loads a single agent from its YAML config.
func (lm *LifecycleManager) LoadAgent(agentID string) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	path := filepath.Join(lm.configDir, agentID+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config %s: %v", path, err)
	}

	var cfg AgentConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse YAML %s: %v", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid agent config %s: %v", path, err)
	}

	agent := BuildAgent(cfg)

	// If already registered, skip or overwrite
	if _, err := lm.manager.Get(agent.ID()); err == nil {
		fmt.Printf("⚠️ Agent '%s' already loaded, overwriting\n", agent.ID())
	}

	if err := lm.manager.Register(agent); err != nil {
		return fmt.Errorf("failed to register agent %s: %v", agent.ID(), err)
	}

	return nil
}

// UnloadAgent removes an agent from the manager.
func (lm *LifecycleManager) UnloadAgent(agentID string) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if _, err := lm.manager.Get(agentID); err != nil {
		return fmt.Errorf("agent %s not found", agentID)
	}

	return lm.manager.Unregister(agentID)
}

// ReloadAgent reloads an agent from disk, replacing the existing instance.
func (lm *LifecycleManager) ReloadAgent(agentID string) error {
	if err := lm.UnloadAgent(agentID); err != nil {
		return fmt.Errorf("failed to unload agent %s: %v", agentID, err)
	}
	return lm.LoadAgent(agentID)
}

// ListAgents returns all currently loaded agent IDs.
func (lm *LifecycleManager) ListAgents() []string {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	agents := lm.manager.List()
	ids := make([]string, len(agents))
	for i, a := range agents {
		ids[i] = a.ID()
	}
	return ids
}

// SaveOrMergeConfig saves the CLI agent config, merging with existing YAML if it exists.
func (lm *LifecycleManager) SaveOrMergeConfig(cfg AgentConfig) error {
	path := filepath.Join(lm.configDir, cfg.ID+".yaml")

	// Load existing config if present
	var existing AgentConfig
	if data, err := os.ReadFile(path); err == nil {
		_ = yaml.Unmarshal(data, &existing)
	}

	// Merge CLI config into existing
	existing.Merge(cfg)

	// Validate
	if err := existing.Validate(); err != nil {
		return err
	}

	// Write back YAML
	yamlData, err := yaml.Marshal(&existing)
	if err != nil {
		return err
	}
	return os.WriteFile(path, yamlData, 0644)
}
