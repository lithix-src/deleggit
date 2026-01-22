package service

import (
	"fmt"
	"os"

	"github.com/datacraft/catalyst/core/internal/adapter/agent"
	"github.com/datacraft/catalyst/core/internal/adapter/workspace"
	"github.com/datacraft/catalyst/core/internal/domain"
	"github.com/point-unknown/catalyst/pkg/mcp"
	"github.com/point-unknown/catalyst/pkg/vector"
	"gopkg.in/yaml.v3"
)

// AgentFactory creates an agent instance from configuration.
func AgentFactory(cfg domain.AgentConfig, llm domain.LLMProvider, registry mcp.Registry, vecStore vector.Store, resolver domain.ContextResolver) (domain.Agent, error) {
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
		// Create Secure Workspace with Resolver (GitGuard)
		ws := workspace.NewLocalWorkspace(cfg.Security, resolver)
		return agent.NewEngineerAgent(cfg.ID, llm, ws, cfg.Security), nil

	case "liaison":
		// Inject Vector Store into Liaison
		return agent.NewLiaisonAgent(cfg.ID, llm, registry, vecStore), nil

	default:
		return nil, fmt.Errorf("unknown agent type: %s", cfg.Type)
	}
}

// LoadAgents reads the YAML config and instantiates agents.
func LoadAgents(path string, llm domain.LLMProvider, registry mcp.Registry, vecStore vector.Store, resolver domain.ContextResolver) ([]domain.Agent, []domain.Mission, error) {
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
		// Use Factory
		a, err := AgentFactory(cfg, llm, registry, vecStore, resolver)
		if err != nil {
			// fallback or skip
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
