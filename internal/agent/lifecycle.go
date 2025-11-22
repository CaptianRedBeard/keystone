// Package agent provides base implementations and helpers for AI agents.
package agent

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"keystone/internal/providers"

	"gopkg.in/yaml.v3"
)

const DefaultAgentsDir = "./agents"

// LifecycleManager manages agent configurations and provider resolution.
type LifecycleManager struct {
	configDir string
	manager   *AgentManager
	providers map[string]providers.Provider
}

// NewLifecycleManager creates a new LifecycleManager with optional config directory and provider map.
func NewLifecycleManager(configDir string, providersMap map[string]providers.Provider) *LifecycleManager {
	if configDir == "" {
		configDir = DefaultAgentsDir
	}
	if providersMap == nil {
		providersMap = make(map[string]providers.Provider)
	}
	return &LifecycleManager{
		configDir: configDir,
		manager:   NewManager(),
		providers: providersMap,
	}
}

// Manager returns the internal AgentManager.
func (lm *LifecycleManager) Manager() *AgentManager {
	return lm.manager
}

// RegisterProvider adds a provider under the specified name.
func (lm *LifecycleManager) RegisterProvider(name string, p providers.Provider) {
	lm.providers[name] = p
}

// ResolveProvider retrieves a registered provider by name.
func (lm *LifecycleManager) ResolveProvider(name string) (providers.Provider, error) {
	p, ok := lm.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return p, nil
}

// SaveOrMergeConfig saves an AgentConfig to YAML in the config directory.
func (lm *LifecycleManager) SaveOrMergeConfig(cfg AgentConfig) error {
	if cfg.ID == "" {
		return fmt.Errorf("agent config must have an ID")
	}
	if err := os.MkdirAll(lm.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	path := filepath.Join(lm.configDir, cfg.ID+".yaml")
	data, _ := yaml.Marshal(&cfg)
	return os.WriteFile(path, data, 0644)
}

// LoadAgent loads a single agent by ID and registers it.
func (lm *LifecycleManager) LoadAgent(agentID string) error {
	if agentID == "" {
		return fmt.Errorf("agentID required")
	}
	path := filepath.Join(lm.configDir, agentID+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("agent config not found: %s", path)
	}

	cfg, err := lm.loadAgentConfig(path)
	if err != nil {
		return fmt.Errorf("failed to load agent config: %w", err)
	}
	provider, err := lm.ResolveProvider(cfg.Provider)
	if err != nil {
		return fmt.Errorf("failed to resolve provider: %w", err)
	}

	agent := NewAgent(cfg.ID, cfg.Name, cfg.Description, provider, cfg.Model, cfg.Memory,
		WithParameters(cfg.Parameters),
		WithPromptTemplate(cfg.PromptTemplate),
		WithLogging(cfg.Logging),
	)
	return lm.manager.Register(agent)
}

// LoadAgentsFromDir loads all YAML agent configs from the config directory.
func (lm *LifecycleManager) LoadAgentsFromDir() error {
	dir := lm.configDir
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("warning: agents directory %q does not exist; continuing", dir)
			return nil
		}
		return fmt.Errorf("failed to read agents directory %q: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("agents path %q is not a directory", dir)
	}

	count := 0
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() || (filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		cfg, err := lm.loadAgentConfig(path)
		if err != nil {
			log.Printf("warning: skipping malformed agent file %s: %v", entry.Name(), err)
			continue
		}
		provider, err := lm.ResolveProvider(cfg.Provider)
		if err != nil {
			log.Printf("warning: skipping agent %s (provider %q not found)", cfg.ID, cfg.Provider)
			continue
		}
		agent := NewAgent(cfg.ID, cfg.Name, cfg.Description, provider, cfg.Model, cfg.Memory,
			WithParameters(cfg.Parameters),
			WithPromptTemplate(cfg.PromptTemplate),
			WithLogging(cfg.Logging),
		)
		if err := lm.manager.Register(agent); err != nil {
			log.Printf("warning: failed to register agent %s: %v", cfg.ID, err)
			continue
		}
		count++
	}
	log.Printf("loaded %d agent(s) from %q", count, dir)
	return nil
}

// loadAgentConfig reads a YAML file and unmarshals it into an AgentConfig.
func (lm *LifecycleManager) loadAgentConfig(path string) (AgentConfig, error) {
	var cfg AgentConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
