package workflow

import (
	"context"
	"testing"

	"keystone/internal/agent"
	"keystone/internal/providers"
	"keystone/internal/tickets"

	"github.com/stretchr/testify/assert"
)

// --- Mock provider implementation ---
type MockProvider struct{}

func (p *MockProvider) GenerateResponse(ctx context.Context, prompt string, model string) (string, error) {
	return "mock response", nil
}

func (p *MockProvider) UsageInfo() (providers.Usage, error) {
	return providers.Usage{}, nil
}

// --- Mock agent implementation ---
type MockAgent struct {
	id      string
	memory  string
	context map[string]string
}

func (a *MockAgent) ID() string          { return a.id }
func (a *MockAgent) Name() string        { return a.id }
func (a *MockAgent) Description() string { return a.id }
func (a *MockAgent) Handle(_ context.Context, input string, t *tickets.Ticket) (string, error) {
	if t != nil {
		for k, v := range a.context {
			t.SetNamespaced(a.id, k, v)
		}
		// Removed IncrementStep to avoid double counting
	}
	return "🧠 Venice says (model=" + a.DefaultModel() + "): \"" + input + "\" [mocked]", nil
}
func (a *MockAgent) Memory() string                 { return a.memory }
func (a *MockAgent) DefaultModel() string           { return "mock-model" }
func (a *MockAgent) Provider() providers.Provider   { return &MockProvider{} }
func (a *MockAgent) PromptTemplate() string         { return "" }
func (a *MockAgent) Parameters() map[string]string  { return map[string]string{} }
func (a *MockAgent) LoggingEnabled() bool           { return false }
func (a *MockAgent) ContextData() map[string]string { return a.context }

// --- Tests ---

func TestWorkflow_RunSequentialAgents(t *testing.T) {
	manager := agent.NewManager()

	agent1 := &MockAgent{id: "a1", context: map[string]string{"foo": "bar"}}
	agent2 := &MockAgent{id: "a2", context: map[string]string{"baz": "qux"}}

	_ = manager.Register(agent1)
	_ = manager.Register(agent2)

	wf := Workflow{
		ID: "wf1",
		Steps: []Step{
			{AgentID: "a1", Input: "input1"},
			{AgentID: "a2", Input: "input2"},
		},
	}

	ticket := tickets.NewTicket("t1", "default", nil)
	engine := NewEngine(manager, true)

	results, err := engine.Run(context.Background(), wf, ticket)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	assert.Equal(t, "a1", results[0].AgentID)
	assert.Contains(t, results[0].Output, "input1")
	assert.Contains(t, results[0].Output, "[mocked]")

	assert.Equal(t, "a2", results[1].AgentID)
	assert.Contains(t, results[1].Output, "input2")
	assert.Contains(t, results[1].Output, "[mocked]")

	ns1, ok1 := ticket.GetNamespaced("a1", "foo")
	ns2, ok2 := ticket.GetNamespaced("a2", "baz")
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, "bar", ns1)
	assert.Equal(t, "qux", ns2)

	// Engine increments ticket once per step
	assert.Equal(t, 2, ticket.Step)
}

func TestWorkflow_RunMissingAgent(t *testing.T) {
	manager := agent.NewManager()
	wf := Workflow{
		ID: "wf2",
		Steps: []Step{
			{AgentID: "missing", Input: "input"},
		},
	}
	ticket := tickets.NewTicket("t2", "default", nil)
	engine := NewEngine(manager, true)

	_, err := engine.Run(context.Background(), wf, ticket)
	assert.Error(t, err)
}

func TestWorkflow_ParamMerging(t *testing.T) {
	manager := agent.NewManager()

	agent1 := &MockAgent{id: "a1", context: map[string]string{}}
	_ = manager.Register(agent1)

	wf := Workflow{
		ID: "wf3",
		Steps: []Step{
			{AgentID: "a1", Input: "input", Params: map[string]string{"param": "override"}},
		},
	}

	ticket := tickets.NewTicket("t3", "default", nil)
	engine := NewEngine(manager, true)

	results, err := engine.Run(context.Background(), wf, ticket)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Contains(t, results[0].Output, "input")
	assert.Contains(t, results[0].Output, "[mocked]")
}
