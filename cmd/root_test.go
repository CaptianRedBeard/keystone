package cmd

import (
	"bytes"
	"strings"
	"testing"

	"keystone/internal/agent"
	"keystone/internal/config"
	"keystone/internal/tickets"
)

func TestRootCommand(t *testing.T) {
	dummyProvider := func(dir string) *agent.AgentManager {
		return agent.NewManager()
	}

	store := tickets.NewStore(t.TempDir())

	runCommand := func(args ...string) string {
		buf := new(bytes.Buffer)
		cfgLoader := func(_ string) (*config.Config, error) { return config.New(), nil }
		cmd := NewRootCmd(dummyProvider, cfgLoader, buf, store)
		cmd.SetArgs(args)
		_ = cmd.Execute()
		return buf.String()
	}

	// Test version output
	output := runCommand("--version")
	if !strings.Contains(output, "v0.1.0") {
		t.Errorf("expected version output, got: %s", output)
	}

	// Test verbose flag
	output = runCommand("--verbose")
	if !strings.Contains(output, "Verbose mode enabled") {
		t.Errorf("expected verbose output, got: %s", output)
	}
}
