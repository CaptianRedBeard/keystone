package cmd

import (
	"bytes"
	"strings"
	"testing"

	"keystone/internal/agent"
	"keystone/internal/config"
)

func TestUsageCLI(t *testing.T) {
	dummyProvider := func(dir string) *agent.AgentManager { return agent.NewManager() }

	runCommand := func(args ...string) string {
		buf := new(bytes.Buffer)
		cfgLoader := func(_ string) (*config.Config, error) { return config.New(), nil }
		cmd := NewRootCmd(dummyProvider, cfgLoader, buf)
		cmd.SetArgs(args)
		_ = cmd.Execute()
		return buf.String()
	}

	// usage summary
	output := runCommand("usage", "summary")
	if !strings.Contains(output, "Requests:") || !strings.Contains(output, "Tokens used:") {
		t.Fatalf("expected output to contain 'Requests:' and 'Tokens used:', got: %s", output)
	}

	// usage summary with --days
	output = runCommand("usage", "summary", "--days", "3")
	if !strings.Contains(output, "Requests:") {
		t.Errorf("expected output to contain 'Requests:' with --days flag, got: %s", output)
	}
}
