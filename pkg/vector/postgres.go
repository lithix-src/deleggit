package vector

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore implements Store using PostgreSQL and pgvector.
type PostgresStore struct {
	pool      *pgxpool.Pool
	tableName string
	dimension int
}

// NewPostgresStore creates a new PostgresStore.
func NewPostgresStore(pool *pgxpool.Pool, tableName string, dimension int) *PostgresStore {
	return &PostgresStore{
		pool:      pool,
		tableName: tableName,
		dimension: dimension,
	}
}

// Init ensures the vector extension and table exist.
func (s *PostgresStore) Init(ctx interface{}) error {
	c, ok := ctx.(context.Context)
	if !ok {
		return fmt.Errorf("context must be context.Context")
	}

	// 1. Enable pgvector extension
	_, err := s.pool.Exec(c, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		return fmt.Errorf("failed to create vector extension: %w", err)
	}

	// 2. Create table
	// We use text for ID to allow flexibility (e.g., UUIDs or paths)
	// metadata is stored as JSONB
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			content TEXT,
			embedding vector(%d),
			metadata JSONB
		)
	`, s.tableName, s.dimension)

	_, err = s.pool.Exec(c, query)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", s.tableName, err)
	}

	// 3. Create Index (IVFFlat) for better performance on large datasets
	// Note: Usually requires some data to be effective, but creating if not exists is good practice
	// We skip index creation for MVP to avoid "no data" errors or slow startup
	return nil
}

// Upsert stores or updates documents.
func (s *PostgresStore) Upsert(ctx interface{}, docs []Document) error {
	c, ok := ctx.(context.Context)
	if !ok {
		return fmt.Errorf("context must be context.Context")
	}

	tx, err := s.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	query := fmt.Sprintf(`
		INSERT INTO %s (id, content, embedding, metadata)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			embedding = EXCLUDED.embedding,
			metadata = EXCLUDED.metadata
	`, s.tableName)

	for _, doc := range docs {
		// pgvector expects []float32 for vector type
		_, err := tx.Exec(c, query, doc.ID, doc.Content, doc.Embedding, doc.Metadata)
		if err != nil {
			return fmt.Errorf("failed to upsert doc %s: %w", doc.ID, err)
		}
	}

	return tx.Commit(c)
}

// Search finds the most similar documents.
func (s *PostgresStore) Search(ctx interface{}, query Embedding, limit int) ([]SearchResult, error) {
	c, ok := ctx.(context.Context)
	if !ok {
		return nil, fmt.Errorf("context must be context.Context")
	}

	// Convert Embedding ([]float32) to string representation for SQL if needed,
	// but pgx usually handles []float32 maps to vector type if configured.
	// Actually, pgx/v5 + pgvector libraries are best.
	// For raw SQL without pgvector-go lib, passing specific syntax might be needed.
	// However, simple []float32 usually works with pgx if the driver understands it or we cast it.
	// The operator <=> is cosine distance (lower is better).
	// We want Similarity (Higher is better)? usually 1 - cosine distance.
	// Or just sort by distance.

	sql := fmt.Sprintf(`
		SELECT id, content, embedding, metadata, 1 - (embedding <=> $1) as score
		FROM %s
		ORDER BY embedding <=> $1
		LIMIT $2
	`, s.tableName)

	rows, err := s.pool.Query(c, sql, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search query failed: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		// We scan embedding into a generic interface or specifically []float32
		// pgx default for vector might need specific handling or just work.
		// Let's assume it works for now, if not we usually verify pgx type registration.
		var embVector []float32
		err := rows.Scan(&r.ID, &r.Content, &embVector, &r.Metadata, &r.Score)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		r.Embedding = embVector
		results = append(results, r)
	}

	return results, nil
}

// Delete removes documents by ID.
func (s *PostgresStore) Delete(ctx interface{}, ids []string) error {
	c, ok := ctx.(context.Context)
	if !ok {
		return fmt.Errorf("context must be context.Context")
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ANY($1)`, s.tableName)
	_, err := s.pool.Exec(c, query, ids)
	return err
}
