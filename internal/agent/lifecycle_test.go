package agent

import (
	"keystone/internal/providers"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLifecycleManager_RegisterResolveProvider(t *testing.T) {
	lm := NewLifecycleManager("", nil)
	mock := &MockProvider{}
	lm.RegisterProvider("mock", mock)

	p, err := lm.ResolveProvider("mock")
	require.NoError(t, err)
	require.Equal(t, mock, p)

	_, err = lm.ResolveProvider("unknown")
	require.Error(t, err)
}

func TestLifecycleManager_LoadAndSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	lm := NewLifecycleManager(tmpDir, map[string]providers.Provider{"mock": &MockProvider{}})

	cfg := AgentConfig{
		ID:       "agent1",
		Name:     "Test Agent",
		Provider: "mock",
		Model:    "model1",
		Memory:   "mem1",
	}

	err := lm.SaveOrMergeConfig(cfg)
	require.NoError(t, err)

	// Load agent
	err = lm.LoadAgent("agent1")
	require.NoError(t, err)

	// Attempt to load non-existent agent
	err = lm.LoadAgent("unknown")
	require.Error(t, err)
}

func TestLifecycleManager_LoadAgentsFromDir(t *testing.T) {
	tmpDir := t.TempDir()
	lm := NewLifecycleManager(tmpDir, map[string]providers.Provider{"mock": &MockProvider{}})

	// Write multiple agent configs
	WriteTempAgentConfig(t, tmpDir, "a1", "mock")
	WriteTempAgentConfig(t, tmpDir, "a2", "mock")

	err := lm.LoadAgentsFromDir()
	require.NoError(t, err)

	agents := lm.Manager().List()
	require.Len(t, agents, 2)

	// Write a broken YAML
	broken := filepath.Join(tmpDir, "broken.yaml")
	_ = os.WriteFile(broken, []byte("{ bad_yaml"), 0644)
	err = lm.LoadAgentsFromDir()
	require.NoError(t, err) // logs a warning but does not fail
}
