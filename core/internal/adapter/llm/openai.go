package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OpenAIAdapter implements domain.LLMProvider for OpenAI-compatible APIs (Ollama, vLLM, etc.)
type OpenAIAdapter struct {
	endpoint    string
	model       string
	apiKey      string
	client      *http.Client
	privateMode bool
}

// Config holds the configuration for the adapter
type Config struct {
	Endpoint    string
	Model       string
	APIKey      string
	PrivateMode bool
}

// Msg represents a chat message
type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RequestPayload represents the OpenAI Chat Completion request
type RequestPayload struct {
	Model    string `json:"model"`
	Messages []Msg  `json:"messages"`
	Stream   bool   `json:"stream"`
}

// ResponsePayload represents the OpenAI Chat Completion response
type ResponsePayload struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenAIAdapter creates a new instance of the adapter
func NewOpenAIAdapter(cfg Config) (*OpenAIAdapter, error) {
	// Validate Private Mode
	if cfg.PrivateMode {
		if !isPrivateIP(cfg.Endpoint) {
			return nil, fmt.Errorf("security alert: private mode is enabled but endpoint %s is public", cfg.Endpoint)
		}
	}

	return &OpenAIAdapter{
		endpoint:    strings.TrimSuffix(cfg.Endpoint, "/"),
		model:       cfg.Model,
		apiKey:      cfg.APIKey,
		privateMode: cfg.PrivateMode,
		client:      &http.Client{Timeout: 60 * time.Second},
	}, nil
}

// GenerateCode sends a prompt to the LLM and returns the response
func (a *OpenAIAdapter) GenerateCode(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", a.endpoint)

	payload := RequestPayload{
		Model: a.model,
		Messages: []Msg{
			{Role: "system", Content: "You are an expert software engineer. Output only code or technical explanations."},
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if a.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ResponsePayload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response from LLM")
	}

	return result.Choices[0].Message.Content, nil
}

// CheckHealth verifies the connection
func (a *OpenAIAdapter) CheckHealth(ctx context.Context) error {
	// Simple models check
	url := fmt.Sprintf("%s/models", a.endpoint) // Standard OAI endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	if a.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: %s", resp.Status)
	}
	return nil
}

// isPrivateIP checks if the hostname in the URL resolves to a private IP
func isPrivateIP(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil {
		return false // Fail safe
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host // No port
	}

	// Localhost check
	if host == "localhost" {
		return true
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() {
			return true
		}
	}
	return false
}
