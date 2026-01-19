package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/datacraft/catalyst/core/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/point-unknown/catalyst/pkg/vector"
)

type PostgresStore struct {
	pool   *pgxpool.Pool
	Vector vector.Store
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

// InitSchema creates the necessary tables if they don't exist
func (s *PostgresStore) InitSchema(ctx context.Context) error {
	log.Println("[STORE] Initializing Schema...")

	// 1. Missions Table (Config)
	queryMissions := `
	CREATE TABLE IF NOT EXISTS missions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		trigger_topic TEXT NOT NULL,
		enabled BOOLEAN DEFAULT TRUE,
		config JSONB NOT NULL DEFAULT '{}'
	);`

	// 2. Event Log Table (History)
	queryEvents := `
	CREATE TABLE IF NOT EXISTS event_log (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		source TEXT NOT NULL,
		type TEXT NOT NULL,
		data JSONB
	);
	CREATE INDEX IF NOT EXISTS idx_event_log_type ON event_log(type);
	CREATE INDEX IF NOT EXISTS idx_event_log_ts ON event_log(timestamp DESC);
	`

	if _, err := s.pool.Exec(ctx, queryMissions); err != nil {
		return fmt.Errorf("failed to create missions table: %w", err)
	}
	if _, err := s.pool.Exec(ctx, queryEvents); err != nil {

		return fmt.Errorf("failed to create event_log table: %w", err)
	}

	// 3. Initialize Vector Store
	// Dimension 768 is standard for nomic-embed-text (Ollama default)
	// We use "memories" as the table name
	vecStore := vector.NewPostgresStore(s.pool, "memories", 768)
	if err := vecStore.Init(ctx); err != nil {
		return fmt.Errorf("failed to init vector store: %w", err)
	}
	s.Vector = vecStore

	log.Println("[STORE] Schema Initialized (including Vector).")
	return nil
}

// SaveEvent persists a CloudEvent to the log
func (s *PostgresStore) SaveEvent(ctx context.Context, event domain.CloudEvent) error {
	query := `INSERT INTO event_log (source, type, data, timestamp) VALUES ($1, $2, $3, $4)`
	// Letting DB generate ID

	_, err := s.pool.Exec(ctx, query, event.Source, event.Type, event.Data, event.Time)
	return err
}
