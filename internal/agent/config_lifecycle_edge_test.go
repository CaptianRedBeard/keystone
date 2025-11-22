// Package agent contains edge tests for LifecycleManager and config handling.
package agent

import (
	"os"
	"path/filepath"
	"testing"

	"keystone/internal/providers"

	"github.com/stretchr/testify/require"
)

func TestLifecycleManager_LoadAgent_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	lm := NewLifecycleManager(tmpDir, map[string]providers.Provider{"mock": &MockProvider{}})

	// Attempt to load agent with empty ID
	err := lm.LoadAgent("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "agentID required")

	// Attempt to load non-existent agent
	err = lm.LoadAgent("does_not_exist")
	require.Error(t, err)
	require.Contains(t, err.Error(), "agent config not found")
}

func TestLifecycleManager_SaveOrMergeConfig_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	lm := NewLifecycleManager(tmpDir, nil)

	// Attempt to save config without ID
	cfg := AgentConfig{
		Name:     "NoIDAgent",
		Provider: "mock",
	}
	err := lm.SaveOrMergeConfig(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must have an ID")
}

func TestLifecycleManager_LoadAgentsFromDir_WithBrokenFiles(t *testing.T) {
	tmpDir := t.TempDir()
	lm := NewLifecycleManager(tmpDir, map[string]providers.Provider{"mock": &MockProvider{}})

	// Valid agent
	WriteTempAgentConfig(t, tmpDir, "valid_agent", "mock")

	// Broken YAML file
	brokenFile := filepath.Join(tmpDir, "broken.yaml")
	_ = os.WriteFile(brokenFile, []byte("{ invalid_yaml"), 0644)

	// Load all agents
	err := lm.LoadAgentsFromDir()
	require.NoError(t, err) // should log warning but not fail

	agents := lm.Manager().List()
	require.Len(t, agents, 1)
	require.Equal(t, "valid_agent", agents[0].ID())
}

func TestLifecycleManager_LoadAgentsFromDir_DuplicateIDs(t *testing.T) {
	tmpDir := t.TempDir()
	lm := NewLifecycleManager(tmpDir, map[string]providers.Provider{"mock": &MockProvider{}})

	// Write two configs with the same ID
	WriteTempAgentConfig(t, tmpDir, "dup_agent", "mock")
	WriteTempAgentConfig(t, tmpDir, "dup_agent", "mock")

	err := lm.LoadAgentsFromDir()
	require.NoError(t, err) // logs a warning for duplicate
	agents := lm.Manager().List()
	require.Len(t, agents, 1) // only one registered
	require.Equal(t, "dup_agent", agents[0].ID())
}

func TestAgentConfig_MergeWithParameters(t *testing.T) {
	a := AgentConfig{ID: "a1", Name: "Name1", Provider: "mock", Parameters: map[string]string{"foo": "bar"}}
	b := AgentConfig{Parameters: map[string]string{"baz": "qux"}}

	a.Merge(b)
	require.Equal(t, "bar", a.Parameters["foo"])
	require.Equal(t, "qux", a.Parameters["baz"])
}

func TestAgentConfig_ValidateDefaults(t *testing.T) {
	cfg := AgentConfig{ID: "id", Name: "n", Provider: "mock"}
	err := cfg.Validate()
	require.NoError(t, err)
	require.Equal(t, "default", cfg.Model)
}
