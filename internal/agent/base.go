// Package agent provides base implementations and helpers for AI agents.
package agent

import (
	"context"
	"fmt"
	"keystone/internal/providers"
	"keystone/internal/tickets"
)

// AgentBase is a simple base implementation of an Agent.
type AgentBase struct {
	id             string
	name           string
	description    string
	memory         string
	model          string
	provider       providers.Provider
	promptTemplate string
	parameters     map[string]string
	logging        bool
}

// AgentOption is a functional option to configure AgentBase.
type AgentOption func(*AgentBase)

// WithPromptTemplate sets the agent's prompt template.
func WithPromptTemplate(tpl string) AgentOption {
	return func(a *AgentBase) { a.promptTemplate = tpl }
}

// WithParameters sets the agent's parameters.
func WithParameters(params map[string]string) AgentOption {
	return func(a *AgentBase) { a.parameters = params }
}

// WithLogging enables or disables agent logging.
func WithLogging(enabled bool) AgentOption {
	return func(a *AgentBase) { a.logging = enabled }
}

// NewAgent creates a new AgentBase with defaults applied.
func NewAgent(id, name, description string, provider providers.Provider, model, memory string, opts ...AgentOption) Agent {
	if id == "" {
		id = "default_agent"
	}
	if name == "" {
		name = "Default Agent"
	}
	if model == "" {
		model = "default"
	}
	if memory == "" {
		memory = "session_sample"
	}

	a := &AgentBase{
		id:          id,
		name:        name,
		description: description,
		provider:    provider,
		model:       model,
		memory:      memory,
		parameters:  make(map[string]string),
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}

// ID returns the agent's ID.
func (a *AgentBase) ID() string { return a.id }

// Name returns the agent's display name.
func (a *AgentBase) Name() string { return a.name }

// Description returns the agent's description.
func (a *AgentBase) Description() string { return a.description }

// Memory returns the agent's memory identifier.
func (a *AgentBase) Memory() string { return a.memory }

// DefaultModel returns the agent's default model.
func (a *AgentBase) DefaultModel() string { return a.model }

// Provider returns the agent's provider instance.
func (a *AgentBase) Provider() providers.Provider { return a.provider }

// PromptTemplate returns the agent's prompt template.
func (a *AgentBase) PromptTemplate() string { return a.promptTemplate }

// Parameters returns the agent's parameters map.
func (a *AgentBase) Parameters() map[string]string { return a.parameters }

// LoggingEnabled returns true if logging is enabled.
func (a *AgentBase) LoggingEnabled() bool { return a.logging }

// Handle processes input using the agent's provider.
func (a *AgentBase) Handle(ctx context.Context, input string, t *tickets.Ticket) (string, error) {
	if a.provider == nil {
		return "", fmt.Errorf("agent %s has no provider configured", a.id)
	}
	return a.provider.GenerateResponse(ctx, input, a.model)
}
