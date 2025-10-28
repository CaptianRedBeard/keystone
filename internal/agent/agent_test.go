package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"keystone/internal/providers"

	"gopkg.in/yaml.v3"
)

// mockProvider implements the Provider interface for testing
type mockProvider struct{}

func (m *mockProvider) GenerateResponse(ctx context.Context, input, model string) (string, error) {
	return "mock response: " + input, nil
}

func (m *mockProvider) UsageInfo() (providers.Usage, error) {
	return providers.Usage{
		Requests: 0,
		Tokens:   0,
	}, nil
}

func TestAgentBase_Handle(t *testing.T) {
	agent := NewAgent("a1", "Test Agent", "desc", &mockProvider{}, "default", "", WithLogging(true))
	resp, err := agent.Handle(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}
	expected := "mock response: hello"
	if resp != expected {
		t.Errorf("Expected '%s', got '%s'", expected, resp)
	}
}

func TestAgentManager_RegisterGetList(t *testing.T) {
	mgr := NewManager()
	agent := NewAgent("a1", "Test Agent", "desc", &mockProvider{}, "default", "")

	if err := mgr.Register(agent); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if _, err := mgr.Get("a1"); err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	agents := mgr.List()
	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}

	// Duplicate registration should fail
	if err := mgr.Register(agent); err == nil {
		t.Errorf("Expected error on duplicate register")
	}
}

func TestLifecycleManager_LoadUnloadReload(t *testing.T) {
	// Create temporary directory for YAML config
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "a1.yaml")

	cfg := AgentConfig{
		ID:          "a1",
		Name:        "Test Agent",
		Description: "Lifecycle test agent",
		Provider:    "venice",
		Model:       "default",
	}

	data, _ := yaml.Marshal(&cfg)
	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test YAML: %v", err)
	}

	lm := NewLifecycleManager(tmpDir)

	// Load agent
	if err := lm.LoadAgent("a1"); err != nil {
		t.Fatalf("LoadAgent failed: %v", err)
	}

	// Reload agent
	if err := lm.ReloadAgent("a1"); err != nil {
		t.Fatalf("ReloadAgent failed: %v", err)
	}

	// Unload agent
	if err := lm.UnloadAgent("a1"); err != nil {
		t.Fatalf("UnloadAgent failed: %v", err)
	}

	// Unload again should fail
	if err := lm.UnloadAgent("a1"); err == nil {
		t.Errorf("Expected error on unloading missing agent")
	}
}

func TestAgentConfig_MergeValidate(t *testing.T) {
	a := AgentConfig{ID: "a1", Name: "A1", Provider: "venice"}
	b := AgentConfig{Name: "A1-updated", Parameters: map[string]string{"x": "1"}, Logging: true}

	a.Merge(b)

	if a.Name != "A1-updated" {
		t.Errorf("Expected name 'A1-updated', got %s", a.Name)
	}
	if !a.Logging {
		t.Errorf("Expected logging to be true")
	}

	if err := a.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Missing ID should fail validation
	a.ID = ""
	if err := a.Validate(); err == nil {
		t.Errorf("Expected validation error for missing ID")
	}
}
