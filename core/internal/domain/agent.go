package domain

import (
	"context"
)

// AgentType defines the execution mode of an agent (The Spectrum)
type AgentType string

const (
	AgentTypeReporter     AgentType = "reporter"     // Level 1: Logs/Telemetry
	AgentTypeCommunicator AgentType = "communicator" // Level 2: Inter-agent signals
	AgentTypeExpressor    AgentType = "expressor"    // Level 3: Artifact generation
)

// Agent allows for the execution of a specific unit of work.
// This is the primary plugin interface for the system.
type Agent interface {
	// ID returns the unique identifier of the agent (e.g., "TrendScout")
	ID() string

	// Type returns the spectrum level of the agent
	Type() AgentType

	// Execute performs the agent's logic given a trigger event.
	// It returns an optional output event or an error.
	Execute(ctx context.Context, input CloudEvent) (*CloudEvent, error)
}

// Mission configuration that maps an input topic to a list of agents.
type Mission struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	TriggerTopic string   `json:"trigger_topic"` // MQTT Topic to subscribe to (e.g. "sensor/cpu/#")
	Agents       []string `json:"agents"`        // List of Agent IDs to execute in order
}
