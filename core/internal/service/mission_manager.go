package service

import (
	"context"
	"fmt"
	"log"

	"github.com/datacraft/catalyst/core/internal/domain"
)

// MissionManager coordinates the execution of Missions based on triggers.
type MissionManager struct {
	registry *AgentRegistry
	missions []domain.Mission
	publish  func(topic string, event domain.CloudEvent)
}

func NewMissionManager(registry *AgentRegistry, publisher func(topic string, event domain.CloudEvent)) *MissionManager {
	return &MissionManager{
		registry: registry,
		missions: make([]domain.Mission, 0),
		publish:  publisher,
	}
}

// LoadMission adds a mission configuration to the active set.
func (m *MissionManager) LoadMission(mission domain.Mission) {
	m.missions = append(m.missions, mission)
	log.Printf("[MISSION] Loaded: %s (Trigger: %s)", mission.Name, mission.TriggerTopic)
}

// ProcessEvent is the main entrypoint for the EventBus.
// It checks if the event matches any Mission triggers and executes the flow.
func (m *MissionManager) ProcessEvent(event domain.CloudEvent) {
	for _, mission := range m.missions {
		// Basic Topic Match (Exact match for now, wildcard support later)
		// Simulating wildcard match for "sensor/#"
		if mission.TriggerTopic == "#" || mission.TriggerTopic == event.Type || (mission.TriggerTopic == "sensor/#" && event.Source != "") {
			m.executeMission(mission, event)
		}
	}
}

func (m *MissionManager) executeMission(mission domain.Mission, trigger domain.CloudEvent) {
	ctx := context.Background()
	log.Printf("[ORCHESTRATOR] Triggering Mission: %s", mission.Name)

	currentPayload := trigger

	for _, agentID := range mission.Agents {
		agent, err := m.registry.Get(agentID)
		if err != nil {
			log.Printf("[ERROR] Mission Failed: %v", err)
			return
		}

		log.Printf("[EXEC] Agent '%s' starting...", agent.ID())
		output, err := agent.Execute(ctx, currentPayload)
		if err != nil {
			log.Printf("[ERROR] Agent '%s' failed: %v", agent.ID(), err)
			return
		}

		// Pipeline: Output becomes input for next agent (if exists)
		if output != nil {
			currentPayload = *output
			log.Printf("[EXEC] Agent '%s' produced output type: %s", agent.ID(), output.Type)

			// Publish the output side-effect
			topic := output.Type // Default fallback
			if output.Type == "tool.call" {
				topic = "tool/call"
			} else {
				topic = fmt.Sprintf("agent/%s/output", agent.ID())
			}
			m.publish(topic, *output)
		}
	}
}
