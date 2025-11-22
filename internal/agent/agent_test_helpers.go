// Package agent provides helper functions and mocks for testing agents.
package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"keystone/internal/providers"
	"keystone/internal/tickets"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// MockProvider is a simple test provider returning predictable responses.
type MockProvider struct{}

func (m *MockProvider) GenerateResponse(ctx context.Context, input, model string) (string, error) {
	return "mock response: " + input, nil
}

func (m *MockProvider) UsageInfo() (providers.Usage, error) {
	return providers.Usage{Requests: 0, Tokens: 0}, nil
}

func (m *MockProvider) Name() string { return "mock" }

// NewMockTicket returns a pre-populated ticket for testing.
func NewMockTicket() *tickets.Ticket {
	return &tickets.Ticket{
		ID:     "t1",
		UserID: "user1",
		Context: map[string]interface{}{
			"test.key": "value",
		},
	}
}

// WriteYAML writes an object to a file as YAML.
func WriteYAML(t *testing.T, path string, data any) {
	t.Helper()
	b, err := yaml.Marshal(data)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, b, 0644))
}

// BuildTestAgent constructs a simple test agent with default memory/model.
func BuildTestAgent(id, name string) Agent {
	return NewAgent(id, name, "Test agent description", &MockProvider{}, id+"-model", "memory-string", WithLogging(true))
}

// HandleInput runs an agent's Handle method and asserts no error occurred.
func HandleInput(t *testing.T, agent Agent, input string) string {
	t.Helper()
	resp, err := agent.Handle(context.Background(), input, NewMockTicket())
	require.NoError(t, err)
	return resp
}

// WriteTempAgentConfig writes a temporary agent config YAML for tests.
func WriteTempAgentConfig(t *testing.T, dir, id, provider string) string {
	t.Helper()
	cfg := AgentConfig{
		ID:          id,
		Name:        id + "-name",
		Description: "temp test agent",
		Provider:    provider,
		Model:       "default",
		Memory:      "none",
	}
	path := filepath.Join(dir, id+".yaml")
	WriteYAML(t, path, cfg)
	return path
}
