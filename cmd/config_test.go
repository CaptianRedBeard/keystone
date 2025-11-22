package cmd

import (
	"bytes"
	"testing"

	"keystone/internal/config"
)

// TestConfigCommand verifies the config command output.
func TestConfigCommand(t *testing.T) {
	mockLoader := func(path string) (*config.Config, error) {
		return &config.Config{
			AgentsDir:  "mock-agents",
			DBPath:     "/tmp/mock.db",
			JSONOutput: false,
			Secrets:    map[string]string{"api_key": "secret"},
		}, nil
	}

	runCommand := func(args ...string) string {
		buf := new(bytes.Buffer)
		cmd := newConfigCmd(mockLoader)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("failed to execute command: %v", err)
		}
		return buf.String()
	}

	output := runCommand()
	if output == "" {
		t.Errorf("expected config output, got empty string")
	}

	output = runCommand("--show-secrets")
	if output == "" {
		t.Errorf("expected config output with secrets, got empty string")
	}
}
