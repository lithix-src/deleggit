package service

import (
	"fmt"
	"sync"

	"github.com/datacraft/catalyst/core/internal/domain"
)

// AgentRegistry manages the lifecycle and lookup of available Agents.
type AgentRegistry struct {
	mu     sync.RWMutex
	agents map[string]domain.Agent
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]domain.Agent),
	}
}

// Register adds an agent to the registry.
func (r *AgentRegistry) Register(agent domain.Agent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[agent.ID()] = agent
}

// Get retrieves an agent by ID.
func (r *AgentRegistry) Get(id string) (domain.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, ok := r.agents[id]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", id)
	}
	return agent, nil
}

// List returns all registered agents.
func (r *AgentRegistry) List() []domain.Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Agent, 0, len(r.agents))
	for _, a := range r.agents {
		list = append(list, a)
	}
	return list
}
