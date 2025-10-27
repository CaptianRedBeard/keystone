package agent

import (
	"context"
	"fmt"

	"keystone/internal/providers/venice"
)

// NewSampleAgent creates a single prototype agent for Phase 1.
func NewSampleAgent() *Agent {
	// Create a mock Venice provider (no real API calls yet)
	provider := venice.New("MOCK_API_KEY", "")

	// Initialize a simple agent using that provider
	agent := New("sample_agent", "Sample Agent", provider, "default")

	return agent
}

// Description provides a short description for CLI listings.
func (a *Agent) Description() string {
	return "A prototype agent using the Venice provider for testing."
}

// Handle is a convenience wrapper for Run.
func (a *Agent) Handle(ctx context.Context, input string) (string, error) {
	return a.Run(ctx, input)
}

// RegisterSampleAgent registers the sample agent with the given manager.
func RegisterSampleAgent(manager *AgentManager) {
	agent := NewSampleAgent()
	if err := manager.Register(agent); err != nil {
		fmt.Printf("⚠️  Failed to register sample agent: %v\n", err)
	} else {
		fmt.Printf("✅ Sample agent registered: %s\n", agent.Name)
	}
}
