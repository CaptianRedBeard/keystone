// internal/agent/interface.go
package agent

import (
	"context"
	"keystone/internal/providers"
)

// Agent defines the methods every agent must implement.
type Agent interface {
	ID() string
	Name() string
	Description() string
	Handle(ctx context.Context, input string) (string, error)
	Memory() string
	DefaultModel() string
	Provider() providers.Provider
	PromptTemplate() string
	Parameters() map[string]string
	LoggingEnabled() bool
}
