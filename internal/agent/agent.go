// internal/agent/agent.go
package agent

import (
	"context"
	"keystone/internal/providers"
)

// AgentBase implements the Agent interface and holds configurable fields.
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

// NewAgent constructs a new AgentBase with functional options.
func NewAgent(
	id, name, description string,
	provider providers.Provider,
	model, memory string,
	options ...AgentOption,
) *AgentBase {
	a := &AgentBase{
		id:          id,
		name:        name,
		description: description,
		provider:    provider,
		model:       model,
		memory:      memory,
		parameters:  make(map[string]string),
	}

	for _, opt := range options {
		opt(a)
	}

	return a
}

// Functional options
type AgentOption func(*AgentBase)

func WithPromptTemplate(tpl string) AgentOption {
	return func(a *AgentBase) {
		a.promptTemplate = tpl
	}
}

func WithParameters(params map[string]string) AgentOption {
	return func(a *AgentBase) {
		a.parameters = params
	}
}

func WithLogging(enabled bool) AgentOption {
	return func(a *AgentBase) {
		a.logging = enabled
	}
}

// Agent interface implementations
func (a *AgentBase) ID() string                    { return a.id }
func (a *AgentBase) Name() string                  { return a.name }
func (a *AgentBase) Description() string           { return a.description }
func (a *AgentBase) Memory() string                { return a.memory }
func (a *AgentBase) DefaultModel() string          { return a.model }
func (a *AgentBase) Provider() providers.Provider  { return a.provider }
func (a *AgentBase) PromptTemplate() string        { return a.promptTemplate }
func (a *AgentBase) Parameters() map[string]string { return a.parameters }
func (a *AgentBase) LoggingEnabled() bool          { return a.logging }

// Handle executes the agent's main function using its provider.
func (a *AgentBase) Handle(ctx context.Context, input string) (string, error) {
	return a.provider.GenerateResponse(ctx, input, a.model)
}
