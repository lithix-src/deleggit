package domain

import (
	"encoding/json"
	"time"
)

// CloudEvent represents the standard event envelope for the Deleggit system.
// It adheres to the CloudEvents JSON format.
type CloudEvent struct {
	ID          string          `json:"id"`
	Source      string          `json:"source"`
	SpecVersion string          `json:"specversion"`
	Type        string          `json:"type"` // e.g., "sensor.cpu.temp", "agent.log"
	Time        time.Time       `json:"time"`
	Data        json.RawMessage `json:"data"` // Polymorphic payload
}

// Helper to create a new event
func NewEvent(source, eventType string, data interface{}) (CloudEvent, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return CloudEvent{}, err
	}
	return CloudEvent{
		Source:      source,
		SpecVersion: "1.0",
		Type:        eventType,
		Time:        time.Now().UTC(),
		Data:        bytes,
	}, nil
}
