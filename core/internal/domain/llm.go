package domain

import "context"

// LLMProvider defines the contract for any Large Language Model service
// (e.g., OpenAI, Ollama, Anthropic).
type LLMProvider interface {
	// GenerateCode sends a prompt and returns the raw text response.
	GenerateCode(ctx context.Context, prompt string) (string, error)

	// CheckHealth verifies the connection to the LLM service.
	CheckHealth(ctx context.Context) error

	// Embed generates a vector embedding for the given text.
	Embed(ctx context.Context, text string) ([]float32, error)
}
