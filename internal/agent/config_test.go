package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAgentConfig_MergeAndValidate(t *testing.T) {
	a := AgentConfig{ID: "a1", Name: "Name1", Provider: "mock"}
	b := AgentConfig{Name: "Name2", Model: "model2", Logging: true}
	a.Merge(b)

	require.Equal(t, "a1", a.ID)
	require.Equal(t, "Name2", a.Name)
	require.Equal(t, "mock", a.Provider)
	require.Equal(t, "model2", a.Model)
	require.True(t, a.Logging)

	// Validate missing required fields
	err := (&AgentConfig{}).Validate()
	require.Error(t, err)

	// Validate defaults model
	cfg := AgentConfig{ID: "id", Name: "n", Provider: "mock"}
	require.NoError(t, cfg.Validate())
	require.Equal(t, "default", cfg.Model)
}
