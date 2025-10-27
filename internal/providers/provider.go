package providers

import "context"

// Provider defines the common interface every AI backend must implement.
type Provider interface {
	GenerateResponse(ctx context.Context, prompt string, model string) (string, error)
	UsageInfo() (Usage, error)
}

type Usage struct {
	Requests int
	Tokens   int
}
