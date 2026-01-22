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

func (s *PostgresStore) Pool() *pgxpool.Pool {
	return s.pool
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

	// 3. Agents Table (Management)
	queryAgents := `
	CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		service TEXT NOT NULL,
		role TEXT NOT NULL,
		config JSONB NOT NULL DEFAULT '{}',
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	if _, err := s.pool.Exec(ctx, queryAgents); err != nil {
		return fmt.Errorf("failed to create agents table: %w", err)
	}

	// 3.1 Seed Agents if empty
	count := 0
	err := s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM agents").Scan(&count)
	if err == nil && count == 0 {
		log.Println("[STORE] Seeding default agents...")
		defaults := []struct {
			ID, Service, Role string
		}{
			{"spec:frontend:01", "Interface", "Frontend Engineering"},
			{"spec:backend:01", "Orchestrator", "Core Systems Engineering"},
			{"spec:infra:01", "Infrastructure", "DevOps & Site Reliability"},
			{"spec:sim:01", "Simulation", "Synthetic Data Generation"},
			{"verif:qa:01", "Compliance", "Quality Assurance & Audit"},
			{"ops:liaison:01", "Liaison", "Human-Swarm Interface"},
			{"eng:pipeline:01", "PipelineArchitect", "CI/CD Engineering"},
			{"arch:system:01", "SystemArchitect", "Council Chair / Governance"},
			{"eng:feature:01", "SoftwareEngineer", "Vibe Coding / Feature Implementation"},
			{"ops:infra:01", "InfrastructureManager", "Target Host Manager"},
		}

		for _, a := range defaults {
			_, err := s.pool.Exec(ctx, `INSERT INTO agents (id, service, role) VALUES ($1, $2, $3)`, a.ID, a.Service, a.Role)
			if err != nil {
				log.Printf("Failed to seed agent %s: %v\n", a.ID, err)
			}
		}
	}

	// 4. ReposTable (Available Repositories)
	queryRepos := `
	CREATE TABLE IF NOT EXISTS repos (
		id TEXT PRIMARY KEY, 
		org TEXT NOT NULL,
		name TEXT NOT NULL,
		default_branch TEXT NOT NULL DEFAULT 'main'
	);`
	if _, err := s.pool.Exec(ctx, queryRepos); err != nil {
		return fmt.Errorf("failed to create repos table: %w", err)
	}

	// 4.1 Seed Repos
	countRepos := 0
	err = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM repos").Scan(&countRepos)
	if err == nil && countRepos == 0 {
		log.Println("[STORE] Seeding default repos...")
		_, err := s.pool.Exec(ctx, `INSERT INTO repos (id, org, name) VALUES ($1, $2, $3)`, "catalyst-core", "catalytic-ai", "catalyst")
		if err != nil {
			log.Printf("Failed to seed repo: %v\n", err)
		}
	}

	// 5. Context Table (Global State Singleton)
	queryContext := `
	CREATE TABLE IF NOT EXISTS context (
		singleton_id BOOL PRIMARY KEY DEFAULT TRUE,
		active_repo_id TEXT REFERENCES repos(id),
		active_branch TEXT NOT NULL DEFAULT 'main',
		CONSTRAINT context_singleton CHECK (singleton_id)
	);`
	if _, err := s.pool.Exec(ctx, queryContext); err != nil {
		return fmt.Errorf("failed to create context table: %w", err)
	}

	// 5.1 Seed Default Context
	_, err = s.pool.Exec(ctx, `INSERT INTO context (singleton_id, active_repo_id, active_branch) VALUES (TRUE, 'catalyst-core', 'main') ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Failed to seed default context: %v\n", err)
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
