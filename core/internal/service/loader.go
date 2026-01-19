package service

import (
	"fmt"
	"os"

	"github.com/datacraft/catalyst/core/internal/adapter/agent"
	"github.com/datacraft/catalyst/core/internal/adapter/workspace"
	"github.com/datacraft/catalyst/core/internal/domain"
	"github.com/point-unknown/catalyst/pkg/mcp"
	"gopkg.in/yaml.v3"
)

// AgentFactory creates an agent instance from configuration.
func AgentFactory(cfg domain.AgentConfig, llm domain.LLMProvider) (domain.Agent, error) {
	switch cfg.Type {
	case "trend-scout":
		threshold := 80.0
		if v, ok := cfg.Config["threshold"].(float64); ok {
			threshold = v
		}
		return agent.NewTrendScout(cfg.ID, threshold), nil

	case "console-reporter":
		return agent.NewConsoleReporter(cfg.ID), nil

	case "engineer":
		// Create Secure Workspace
		ws := workspace.NewLocalWorkspace(cfg.Security)
		return agent.NewEngineerAgent(cfg.ID, llm, ws, cfg.Security), nil

	default:
		return nil, fmt.Errorf("unknown agent type: %s", cfg.Type)
	}
}

// LoadAgents reads the YAML config and instantiates agents.
func LoadAgents(path string, llm domain.LLMProvider, registry mcp.Registry) ([]domain.Agent, []domain.Mission, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var sysCfg domain.SystemConfig
	if err := yaml.Unmarshal(data, &sysCfg); err != nil {
		return nil, nil, err
	}

	var agents []domain.Agent
	for _, cfg := range sysCfg.Agents {
		var a domain.Agent
		switch cfg.Type {
		case "engineer":
			// We can pass a real workspace here later
			ws := workspace.NewLocalWorkspace(cfg.Security)
			a = agent.NewEngineerAgent(cfg.ID, llm, ws, cfg.Security)
		case "liaison":
			a = agent.NewLiaisonAgent(cfg.ID, llm, registry)
		default:
			// Fallback or generic agent
			continue
		}
		agents = append(agents, a)
	}

	var missions []domain.Mission
	for _, mCfg := range sysCfg.Missions {
		missions = append(missions, domain.Mission{
			ID:           mCfg.ID,
			Name:         mCfg.Name,
			TriggerTopic: mCfg.TriggerTopic,
			Agents:       mCfg.Agents,
		})
	}

	return agents, missions, nil
}
