package domain

// SecurityConfig defines the restrictions for an Agent's execution environment.
type SecurityConfig struct {
	SandboxPath     string   `yaml:"sandbox_path"`     // Absolute path agent is confined to
	AllowedPatterns []string `yaml:"allowed_patterns"` // e.g. ["*.go", "README.md"]
	DeniedPatterns  []string `yaml:"denied_patterns"`  // e.g. [".env", "id_rsa"]
	ReadOnly        bool     `yaml:"read_only"`
}

// Workspace defines the interface for safe file I/O operations.
type Workspace interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	List(path string) ([]string, error)
}
