// Package agent contains tests for AgentBase.
package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAgentBase_GettersAndHandle verifies getters and Handle method work correctly.
func TestAgentBase_GettersAndHandle(t *testing.T) {
	agent := BuildTestAgent("a1", "Agent One")

	// Getter tests
	require.Equal(t, "a1", agent.ID())
	require.Equal(t, "Agent One", agent.Name())
	require.Equal(t, "Test agent description", agent.Description())
	require.Equal(t, "memory-string", agent.Memory())
	require.Equal(t, "a1-model", agent.DefaultModel())
	require.True(t, agent.LoggingEnabled())

	// Handle test
	resp, err := agent.Handle(context.Background(), "hello", NewMockTicket())
	require.NoError(t, err)
	require.Equal(t, "mock response: hello", resp)
}

// TestAgentBase_Handle_NoProvider verifies that Handle returns error when provider is nil.
func TestAgentBase_Handle_NoProvider(t *testing.T) {
	agent := NewAgent("a2", "NoProvider", "desc", nil, "a2-model", "mem")
	_, err := agent.Handle(context.Background(), "input", NewMockTicket())
	require.Error(t, err)
}
