package agent

import (
	"context"
	"fmt"
	"time"

	"keystone/internal/providers"
)

type Agent struct {
	ID       string             // unique identifier
	Name     string             // friendly name
	Provider providers.Provider // the AI provider client
	Model    string             // model or endpoint
	Created  time.Time          // timestamp
}

func New(id, name string, provider providers.Provider, model string) *Agent {
	return &Agent{
		ID:       id,
		Name:     name,
		Provider: provider,
		Model:    model,
		Created:  time.Now(),
	}
}

func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
	if a.Provider == nil {
		return "", fmt.Errorf("agent %s has no provider configured", a.Name)
	}

	resp, err := a.Provider.GenerateResponse(ctx, prompt, a.Model)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return resp, nil
}
