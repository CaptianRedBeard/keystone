package agent

import (
	"context"
	"testing"
)

func TestRegisterAndGetAgent(t *testing.T) {
	manager := NewManager()
	RegisterSampleAgent(manager)

	// Retrieve the agent
	a, err := manager.Get("sample_agent")
	if err != nil {
		t.Fatalf("expected sample_agent, got error: %v", err)
	}

	// Validate properties
	if a.Name != "sample_agent" {
		t.Errorf("expected Name 'sample_agent', got %s", a.Name)
	}
	if desc := a.Description(); desc == "" {
		t.Error("expected non-empty Description")
	}

	// Test Run and Handle
	resp, err := a.Run(context.Background(), "Hello test")
	if err != nil {
		t.Errorf("Run returned error: %v", err)
	}
	if resp == "" {
		t.Error("expected non-empty response")
	}

	resp2, err := a.Handle(context.Background(), "Hello test 2")
	if err != nil {
		t.Errorf("Handle returned error: %v", err)
	}
	if resp2 == "" {
		t.Error("expected non-empty response from Handle")
	}
}

func TestDuplicateRegistration(t *testing.T) {
	manager := NewManager()
	RegisterSampleAgent(manager)
	err := manager.Register(NewSampleAgent())
	if err == nil {
		t.Error("expected error when registering duplicate agent, got nil")
	}
}

func TestGetMissingAgent(t *testing.T) {
	manager := NewManager()
	_, err := manager.Get("nonexistent")
	if err == nil {
		t.Error("expected error when retrieving missing agent, got nil")
	}
}
