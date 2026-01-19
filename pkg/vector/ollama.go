package vector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// OllamaProvider implements Provider using a local Ollama instance.
type OllamaProvider struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

// NewOllamaProvider creates a new provider. url defaults to "http://localhost:11434".
func NewOllamaProvider(url string, model string) *OllamaProvider {
	if url == "" {
		url = "http://localhost:11434"
	}
	if model == "" {
		model = "nomic-embed-text" // Standard open source embedding model
	}
	return &OllamaProvider{
		BaseURL: url,
		Model:   model,
		Client:  &http.Client{},
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaResponse struct {
	Embedding []float64 `json:"embedding"` // Ollama returns float64
}

// Embed generates an embedding for a single text.
func (p *OllamaProvider) Embed(text string) (Embedding, error) {
	reqBody := ollamaRequest{
		Model:  p.Model,
		Prompt: text,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := p.Client.Post(p.BaseURL+"/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status: %s", resp.Status)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Convert float64 to float32
	embedding := make(Embedding, len(result.Embedding))
	for i, v := range result.Embedding {
		embedding[i] = float32(v)
	}

	return embedding, nil
}

// EmbedBatch generates embeddings for multiple texts.
// Ollama doesn't support batch embeddings natively (one call per text usually),
// so we iterate.
func (p *OllamaProvider) EmbedBatch(texts []string) ([]Embedding, error) {
	var embeddings []Embedding
	for _, text := range texts {
		emb, err := p.Embed(text)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, emb)
	}
	return embeddings, nil
}
