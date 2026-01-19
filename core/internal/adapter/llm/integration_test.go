package llm_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/datacraft/catalyst/core/internal/adapter/llm"
)

func TestOpenAIAdapter_Integration(t *testing.T) {
	// Skip if not running in integration mode
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("Skipping integration test (set TEST_INTEGRATION=true)")
	}

	cfg := llm.Config{
		Endpoint:    "http://localhost:11434/v1",
		Model:       "qwen2.5-coder:7b-instruct",
		PrivateMode: true,
	}

	// Allow override
	if ep := os.Getenv("LLM_ENDPOINT"); ep != "" {
		cfg.Endpoint = ep
	}
	if md := os.Getenv("LLM_MODEL"); md != "" {
		cfg.Model = md
	}

	adapter, err := llm.NewOpenAIAdapter(cfg)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Health Check
	if err := adapter.CheckHealth(ctx); err != nil {
		t.Fatalf("Health check failed (is Ollama running?): %v", err)
	}
	t.Log("✅ Health check passed")

	// 2. Generation Check
	prompt := "Write a Go function that adds two integers."
	response, err := adapter.GenerateCode(ctx, prompt)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if len(response) == 0 {
		t.Fatal("Received empty response from LLM")
	}

	t.Logf("✅ LLM Response:\n%s", response)
}
