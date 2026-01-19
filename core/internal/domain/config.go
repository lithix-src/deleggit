package domain

// AgentConfig represents the configuration for a single agent instance.
type AgentConfig struct {
	ID       string                 `yaml:"id"`
	Type     string                 `yaml:"type"`   // e.g., "trend-scout", "engineer"
	Config   map[string]interface{} `yaml:"config"` // Agent-specific settings
	Security SecurityConfig         `yaml:"safety"` // Mandatory Safety Protocol
}

// MissionConfig represents the configuration for a mission.
type MissionConfig struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	TriggerTopic string   `yaml:"trigger_topic"`
	Agents       []string `yaml:"agents"`
}

// SystemConfig is the top-level structure for config/agents.yaml
type SystemConfig struct {
	Agents   []AgentConfig   `yaml:"agents"`
	Missions []MissionConfig `yaml:"missions"`
}
