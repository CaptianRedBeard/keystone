package agent

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"keystone/internal/logger"
	"keystone/internal/providers"
	"keystone/internal/providers/venice"

	"gopkg.in/yaml.v3"
)

// AgentConfig defines the structure of agent YAML files.
type AgentConfig struct {
	ID             string            `yaml:"id"`
	Name           string            `yaml:"name"`
	Description    string            `yaml:"description"`
	Provider       string            `yaml:"provider"`
	Model          string            `yaml:"model"`
	Memory         string            `yaml:"memory"`
	PromptTemplate string            `yaml:"prompt_template,omitempty"`
	Parameters     map[string]string `yaml:"parameters,omitempty"`
	Logging        bool              `yaml:"logging,omitempty"`
}

// ProviderFactory defines a function that returns a Provider.
type ProviderFactory func() providers.Provider

// providerRegistry maps provider names to constructors.
var providerRegistry = map[string]ProviderFactory{
	"venice": func() providers.Provider {
		return venice.New("MOCK_API_KEY", "")
	},
}

// LoadAgentsFromConfig scans a directory and loads all YAML agents.
// This keeps the old signature so cmd files continue to work.
func LoadAgentsFromConfig(manager *AgentManager, configDir string) error {
	return filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("error reading agent config %s: %v", path, readErr)
		}

		var cfg AgentConfig
		if yamlErr := yaml.Unmarshal(data, &cfg); yamlErr != nil {
			return fmt.Errorf("error parsing YAML for %s: %v", path, yamlErr)
		}

		// Skip duplicate IDs
		if _, err := manager.Get(cfg.ID); err == nil {
			fmt.Printf("⚠️ Skipping duplicate agent ID: %s\n", cfg.ID)
			return nil
		}

		// Create provider
		providerFactory, ok := providerRegistry[cfg.Provider]
		if !ok {
			return fmt.Errorf("unknown provider '%s' for agent %s", cfg.Provider, cfg.ID)
		}
		provider := providerFactory()

		base := NewAgent(
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

		if err := manager.Register(base); err != nil {
			return fmt.Errorf("failed to register agent %s: %v", cfg.ID, err)
		}

		logMsg := fmt.Sprintf("Loaded agent: %-20s (%s)", cfg.Name, cfg.ID)
		logger.Log(cfg.ID, logMsg)
		fmt.Println(logMsg)

		return nil
	})
}

// Merge merges another AgentConfig (src) into this one, prioritizing non-empty fields from src.
func (dst *AgentConfig) Merge(src AgentConfig) {
	if src.ID != "" {
		dst.ID = src.ID
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
	if src.Description != "" {
		dst.Description = src.Description
	}
	if src.Provider != "" {
		dst.Provider = src.Provider
	}
	if src.Model != "" {
		dst.Model = src.Model
	}
	if src.Memory != "" {
		dst.Memory = src.Memory
	}
	if src.PromptTemplate != "" {
		dst.PromptTemplate = src.PromptTemplate
	}
	if src.Parameters != nil {
		if dst.Parameters == nil {
			dst.Parameters = make(map[string]string)
		}
		for k, v := range src.Parameters {
			dst.Parameters[k] = v
		}
	}
	if src.Logging {
		dst.Logging = true
	}
}

// Validate ensures required fields are set correctly.
func (cfg *AgentConfig) Validate() error {
	if cfg.ID == "" {
		return fmt.Errorf("missing agent ID")
	}
	if cfg.Name == "" {
		return fmt.Errorf("missing agent name")
	}
	if cfg.Provider == "" {
		return fmt.Errorf("missing provider for agent %s", cfg.ID)
	}
	if cfg.Model == "" {
		cfg.Model = "default"
	}
	return nil
}
