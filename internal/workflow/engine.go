package workflow

import (
	"context"
	"fmt"
	"strings"

	"keystone/internal/agent"
	"keystone/internal/logger"
	"keystone/internal/tickets"
)

// Engine coordinates workflow execution
type Engine struct {
	manager *agent.AgentManager
	verbose bool
}

// NewEngine creates a new workflow engine
func NewEngine(manager *agent.AgentManager, verbose bool) *Engine {
	return &Engine{manager: manager, verbose: verbose}
}

// Run executes a workflow sequentially, updating the ticket after each step
func (e *Engine) Run(ctx context.Context, wf Workflow, ticket *tickets.Ticket) ([]StepResult, error) {
	results := make([]StepResult, 0, len(wf.Steps))
	var prevOutput string

	logger.Info(fmt.Sprintf("Starting workflow '%s' with %d steps", wf.ID, len(wf.Steps)), false)

	for i, step := range wf.Steps {
		a, err := e.manager.Get(step.AgentID)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get agent '%s': %v", step.AgentID, err), false)
			return results, fmt.Errorf("failed to get agent %s: %w", step.AgentID, err)
		}

		// Merge step params with agent default params
		params := make(map[string]string)
		for k, v := range a.Parameters() {
			params[k] = v
		}
		for k, v := range step.Params {
			params[k] = v
		}

		// Determine final input
		finalInput := step.Input
		if finalInput == "" && i > 0 {
			finalInput = prevOutput // default to previous step output
		}

		// Apply agent prompt template
		if a.PromptTemplate() != "" {
			template := a.PromptTemplate()
			for k, v := range params {
				template = replacePlaceholder(template, k, v)
			}
			finalInput = fmt.Sprintf("%s\n%s", template, finalInput)
		}

		logger.Info(fmt.Sprintf("Running step %d - Agent '%s'", i, a.ID()), false)

		// Run the agent
		output, err := a.Handle(ctx, finalInput, ticket)
		if err != nil {
			logger.Error(fmt.Sprintf("Agent '%s' failed: %v", a.ID(), err), false)
			results = append(results, StepResult{AgentID: a.ID(), Output: "", Error: err})
			return results, fmt.Errorf("agent %s failed: %w", a.ID(), err)
		}

		// Update ticket namespace for agent context
		if ca, ok := a.(agent.ContextualAgent); ok {
			for k, v := range ca.ContextData() {
				ticket.SetNamespaced(a.ID(), k, v)
			}
		}

		// Increment step, pass verbose flag from Engine
		ticket.IncrementStep(e.verbose)
		results = append(results, StepResult{AgentID: a.ID(), Output: output, Error: nil})

		logger.Info(fmt.Sprintf("Step %d - Agent '%s' output:\n%s", i, a.ID(), output), false)

		// Save output for chaining
		prevOutput = output
	}

	logger.Info(fmt.Sprintf("Workflow '%s' completed successfully", wf.ID), false)
	return results, nil
}

// Helper: replace {{key}} in template with value
func replacePlaceholder(template, key, val string) string {
	return strings.ReplaceAll(template, fmt.Sprintf("{{%s}}", key), val)
}
