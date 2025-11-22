// Package agent provides base implementations and helpers for AI agents.
package agent

import "fmt"

// AgentConfig defines the structure of an agent YAML configuration.
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

// Validate ensures required fields are set correctly and applies defaults.
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
