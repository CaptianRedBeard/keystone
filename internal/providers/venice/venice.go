package venice

import (
	"context"
	"fmt"

	"keystone/internal/providers"
)

type VeniceProvider struct {
	apiKey  string
	baseURL string
	usage   providers.Usage
}

func New(apiKey, baseURL string) *VeniceProvider {
	return &VeniceProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		usage:   providers.Usage{},
	}
}

// GenerateResponse simulates an API call to the Venice service.
// For Phase 1, this is mocked to return a test response.
func (v *VeniceProvider) GenerateResponse(ctx context.Context, prompt string, model string) (string, error) {
	v.usage.Requests++
	v.usage.Tokens += len(prompt) / 4 // crude token estimate

	response := fmt.Sprintf("🧠 Venice says (model=%s): %q [mocked]", model, prompt)
	return response, nil
}

// UsageInfo returns mock usage data.
func (v *VeniceProvider) UsageInfo() (providers.Usage, error) {
	return v.usage, nil
}
