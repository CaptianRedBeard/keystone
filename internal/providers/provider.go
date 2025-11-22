package providers

import "context"

type Provider interface {
	GenerateResponse(ctx context.Context, prompt string, model string) (string, error)
	UsageInfo() (Usage, error)
}

type Usage struct {
	Requests int
	Tokens   int
}
