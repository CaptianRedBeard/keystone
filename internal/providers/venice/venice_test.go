package venice

import (
	"context"
	"testing"
)

func TestGenerateResponseAndUsage(t *testing.T) {
	provider := New("MOCK_KEY", "")

	// Basic call
	resp, err := provider.GenerateResponse(context.Background(), "Hello world", "default")
	if err != nil {
		t.Fatalf("GenerateResponse error: %v", err)
	}
	if resp == "" {
		t.Fatal("expected non-empty response")
	}

	usage, err := provider.UsageInfo()
	if err != nil {
		t.Fatalf("UsageInfo error: %v", err)
	}

	if usage.Requests != 1 {
		t.Errorf("expected 1 request, got %d", usage.Requests)
	}
	if usage.Tokens == 0 {
		t.Error("expected tokens > 0")
	}

	// Empty prompt
	resp, err = provider.GenerateResponse(context.Background(), "", "default")
	if err != nil {
		t.Errorf("GenerateResponse failed on empty input: %v", err)
	}
	if resp == "" {
		t.Error("expected non-empty response even for empty input")
	}

	// Multiple calls increment usage
	_, _ = provider.GenerateResponse(context.Background(), "Another prompt", "default")
	usage, _ = provider.UsageInfo()
	if usage.Requests != 3 {
		t.Errorf("expected 3 requests after 3 calls, got %d", usage.Requests)
	}
}
