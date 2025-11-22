// Package agent provides base implementations and helpers for AI agents.
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

// BuildAgent constructs an Agent from an AgentConfig.
func BuildAgent(cfg AgentConfig) Agent {
	var provider providers.Provider
	switch cfg.Provider {
	case "venice":
		provider = venice.New("MOCK_API_KEY", "")
	default:
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

// LoadAgentsFromConfig scans a directory and loads all YAML agent configs into the manager.
func LoadAgentsFromConfig(manager *AgentManager, configDir string) error {
	var loadErrs []error

	if _, err := os.Stat(configDir); err != nil {
		if os.IsNotExist(err) {
			logger.Warn(fmt.Sprintf("Agent config directory does not exist: %s", configDir), false)
			return nil
		}
		return fmt.Errorf("failed to access agent config dir %s: %w", configDir, err)
	}

	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			logger.Error(fmt.Sprintf("Error reading agent config %s: %v", path, readErr), false)
			loadErrs = append(loadErrs, fmt.Errorf("error reading agent config %s: %w", path, readErr))
			return nil
		}

		var cfg AgentConfig
		if yamlErr := yaml.Unmarshal(data, &cfg); yamlErr != nil {
			logger.Error(fmt.Sprintf("Error parsing YAML for %s: %v", path, yamlErr), false)
			loadErrs = append(loadErrs, fmt.Errorf("error parsing YAML for %s: %w", path, yamlErr))
			return nil
		}

		if err := cfg.Validate(); err != nil {
			logger.Error(fmt.Sprintf("Invalid agent config %s: %v", path, err), false)
			loadErrs = append(loadErrs, fmt.Errorf("invalid agent config %s: %w", path, err))
			return nil
		}

		if _, err := manager.Get(cfg.ID); err == nil {
			logger.Warn(fmt.Sprintf("Skipping duplicate agent ID: %s", cfg.ID), false)
			return nil
		}

		a := BuildAgent(cfg)
		if err := manager.Register(a); err != nil {
			logger.Error(fmt.Sprintf("Failed to register agent %s: %v", cfg.ID, err), false)
			loadErrs = append(loadErrs, fmt.Errorf("failed to register agent %s: %w", cfg.ID, err))
			return nil
		}

		logger.Info(fmt.Sprintf("Loaded agent: %-20s (%s)", cfg.Name, cfg.ID), false)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking agent config dir: %w", err)
	}
	if len(loadErrs) > 0 {
		return fmt.Errorf("errors occurred while loading agents: %v", loadErrs)
	}
	return nil
}

// LoadDefaultAgent ensures a fallback dummy agent exists in the manager.
func LoadDefaultAgent(manager *AgentManager) {
	const dummyID = "dummy"
	if _, err := manager.Get(dummyID); err == nil {
		return // already exists
	}

	cfg := AgentConfig{
		ID:             dummyID,
		Name:           "Dummy Agent",
		Description:    "A fallback agent for testing and defaults",
		Provider:       "venice",
		Model:          "default",
		PromptTemplate: "{{input}}",
		Parameters:     map[string]string{},
	}
	a := BuildAgent(cfg)
	if err := manager.Register(a); err != nil {
		logger.Warn(fmt.Sprintf("Failed to register default dummy agent: %v", err), false)
		return
	}

	logger.Info("Loaded default dummy agent", false)
}
