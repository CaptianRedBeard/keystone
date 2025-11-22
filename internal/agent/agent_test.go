// Package agent contains core agent tests.
package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAgent_Handle verifies agent handles input and returns a response.
func TestAgent_Handle(t *testing.T) {
	agent := BuildTestAgent("default_agent", "Default Agent")
	resp := HandleInput(t, agent, "hello world")
	require.Equal(t, "mock response: hello world", resp)
}

// TestAgent_Getters verifies all getters return expected values.
func TestAgent_Getters(t *testing.T) {
	a := BuildTestAgent("default_agent", "Default Agent")
	require.Equal(t, "default_agent", a.ID())
	require.Equal(t, "Default Agent", a.Name())
	require.Equal(t, "memory-string", a.Memory())
	require.Equal(t, "default_agent-model", a.DefaultModel())
	require.True(t, a.LoggingEnabled())
	require.Equal(t, map[string]string{}, a.Parameters())
	require.Equal(t, "", a.PromptTemplate())
}

// TestAgent_HandleWithNilProvider ensures agent returns error when provider is nil.
func TestAgent_HandleWithNilProvider(t *testing.T) {
	agent := NewAgent("no_provider", "No Provider Agent", "Agent with nil provider", nil, "no_provider-model", "mem")
	resp, err := agent.Handle(context.Background(), "input", NewMockTicket())
	require.Error(t, err)
	require.Contains(t, err.Error(), "no provider")
	require.Empty(t, resp)
}
