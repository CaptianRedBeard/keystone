package cmd

import (
	"bytes"
	"strings"
	"testing"

	"keystone/internal/agent"
)

// containsError returns true if CLI output contains an error string
func containsError(output string) bool {
	return strings.Contains(output, "Error") || strings.Contains(output, "failed")
}

// mockManagerProvider returns a manager with a single mock agent registered
func mockManagerProvider() *agent.AgentManager {
	manager := agent.NewManager()

	mockAgent := agent.NewAgent(
		"mock1",
		"Mock Agent",
		"A test agent for workflow simulation",
		nil,
		"default-model",
		"default-memory",
	)

	manager.Register(mockAgent)
	return manager
}

// runWorkflowCLI executes workflow commands through the CLI and captures output
func runWorkflowCLI(t *testing.T, args ...string) string {
	t.Helper()
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(func(dir string) *agent.AgentManager {
		return mockManagerProvider()
	}, nil, buf) // nil loader uses defaults for tests

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("workflow command failed: %v", err)
	}
	return buf.String()
}
