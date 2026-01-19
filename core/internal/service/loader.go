package service

import (
	"fmt"
	"os"

	"github.com/datacraft/catalyst/core/internal/adapter/agent"
	"github.com/datacraft/catalyst/core/internal/adapter/workspace"
	"github.com/datacraft/catalyst/core/internal/domain"
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

// LoadAgents reads the YAML configuration and returns a list of Agents.
func LoadAgents(path string, llm domain.LLMProvider) ([]domain.Agent, []domain.Mission, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var sysConfig domain.SystemConfig
	if err := yaml.Unmarshal(data, &sysConfig); err != nil {
		return nil, nil, err
	}

	var agents []domain.Agent
	for _, agentCfg := range sysConfig.Agents {
		a, err := AgentFactory(agentCfg, llm)
		if err != nil {
			return nil, nil, err
		}
		agents = append(agents, a)
	}

	var missions []domain.Mission
	for _, mCfg := range sysConfig.Missions {
		missions = append(missions, domain.Mission{
			ID:           mCfg.ID,
			Name:         mCfg.Name,
			TriggerTopic: mCfg.TriggerTopic,
			Agents:       mCfg.Agents,
		})
	}

	return agents, missions, nil
}
