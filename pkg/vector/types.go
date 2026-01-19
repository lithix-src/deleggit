package vector

// Embedding represents a vector of floats.
type Embedding []float32

// Document represents a chunk of text with its vector embedding and metadata.
type Document struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Embedding Embedding              `json:"embedding,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// SearchResult represents a document returned from a similarity search.
type SearchResult struct {
	Document
	Score float32 `json:"score"`
}

// Provider defines the interface for generating embeddings.
type Provider interface {
	Embed(text string) (Embedding, error)
	EmbedBatch(texts []string) ([]Embedding, error)
}

// Store defines the interface for storing and retrieving vectors.
type Store interface {
	// Init ensures the vector store is ready (e.g., creates extensions/tables).
	Init(ctx interface{}) error
	// Upsert stores or updates documents.
	Upsert(ctx interface{}, docs []Document) error
	// Search finds the most similar documents to the query vector.
	Search(ctx interface{}, query Embedding, limit int) ([]SearchResult, error)
	// Delete removes documents by ID.
	Delete(ctx interface{}, ids []string) error
}
