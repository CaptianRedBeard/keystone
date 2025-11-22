package agent

import (
	"context"
	"keystone/internal/providers"
	"keystone/internal/tickets"
)

type Agent interface {
	ID() string
	Name() string
	Description() string
	Handle(ctx context.Context, input string, t *tickets.Ticket) (string, error)
	Memory() string
	DefaultModel() string
	Provider() providers.Provider
	PromptTemplate() string
	Parameters() map[string]string
	LoggingEnabled() bool
}

type ContextualAgent interface {
	Agent
	ContextData() map[string]string
}
