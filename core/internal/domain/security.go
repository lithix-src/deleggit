package domain

import (
	"context"
)

// SecurityConfig defines the restrictions for an Agent's execution environment.
type SecurityConfig struct {
	SandboxPath     string   `yaml:"sandbox_path"`     // Absolute path agent is confined to
	AllowedPatterns []string `yaml:"allowed_patterns"` // e.g. ["*.go", "README.md"]
	DeniedPatterns  []string `yaml:"denied_patterns"`  // e.g. [".env", "id_rsa"]
	ReadOnly        bool     `yaml:"read_only"`
}

// ContextResolver defines how to retrieve the dynamic workspace root
type ContextResolver interface {
	GetActiveRepoPath(ctx context.Context) (string, error)
}

// Workspace defines the interface for safe file I/O operations.
type Workspace interface {
	ReadFile(ctx context.Context, path string) ([]byte, error)
	WriteFile(ctx context.Context, path string, data []byte) error
	List(ctx context.Context, path string) ([]string, error)
}
