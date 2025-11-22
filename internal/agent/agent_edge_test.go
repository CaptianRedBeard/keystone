// Package agent contains tests for agent behaviors, including edge cases.
package agent

import (
	"context"
	"errors"
	"testing"

	"keystone/internal/providers"

	"github.com/stretchr/testify/require"
)

// ErrProviderFail is used to simulate a failing provider.
var ErrProviderFail = errors.New("provider failure")

// FailingProvider simulates a provider that always fails.
type FailingProvider struct{}

func (f *FailingProvider) GenerateResponse(ctx context.Context, input, model string) (string, error) {
	return "", ErrProviderFail
}

func (f *FailingProvider) UsageInfo() (providers.Usage, error) {
	return providers.Usage{}, nil
}

func (f *FailingProvider) Name() string { return "failing" }

// TestAgent_HandleWithProviderError ensures the agent propagates provider errors.
func TestAgent_HandleWithProviderError(t *testing.T) {
	agent := NewAgent("err1", "FailingAgent", "desc", &FailingProvider{}, "model", "mem")
	_, err := agent.Handle(context.Background(), "input", NewMockTicket())
	require.ErrorIs(t, err, ErrProviderFail)
}

// TestAgent_HandleNilProvider ensures handling an agent with nil provider returns error.
func TestAgent_HandleNilProvider(t *testing.T) {
	agent := NewAgent("nil1", "NilProviderAgent", "desc", nil, "model", "mem")
	_, err := agent.Handle(context.Background(), "input", NewMockTicket())
	require.Error(t, err)
	require.Contains(t, err.Error(), "no provider")
}

// TestAgent_EmptyIDNameDefaults ensures default values are set for empty ID/Name.
func TestAgent_EmptyIDNameDefaults(t *testing.T) {
	a := BuildTestAgent("", "")
	require.Equal(t, "default_agent", a.ID())
	require.Equal(t, "Default Agent", a.Name())
}

// TestAgent_ParameterMerging ensures options correctly override parameters.
func TestAgent_ParameterMerging(t *testing.T) {
	a := BuildTestAgent("p1", "ParamAgent")
	WithParameters(map[string]string{"foo": "bar"})(a.(*AgentBase))
	require.Equal(t, "bar", a.Parameters()["foo"])
}

// TestAgent_PromptTemplateOption ensures prompt template option is applied.
func TestAgent_PromptTemplateOption(t *testing.T) {
	a := BuildTestAgent("tpl1", "TemplateAgent")
	WithPromptTemplate("Hello {{.Input}}")(a.(*AgentBase))
	require.Equal(t, "Hello {{.Input}}", a.PromptTemplate())
}
