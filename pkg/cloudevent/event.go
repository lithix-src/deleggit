package cloudevent

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event represents the Catalyst Service Standard (CSS) data protocol.
// It adheres to the CloudEvent JSON format.
type Event struct {
	SpecVersion string          `json:"specversion"`
	ID          string          `json:"id"`
	Source      string          `json:"source"`
	Type        string          `json:"type"`
	Time        time.Time       `json:"time"`
	Data        json.RawMessage `json:"data,omitempty"`
}

// New creates a new Event with standard defaults (ID, Time, SpecVersion).
func New(source, eventType string, data interface{}) (Event, error) {
	rawBytes, err := json.Marshal(data)
	if err != nil {
		return Event{}, err
	}

	return Event{
		SpecVersion: "1.0",
		ID:          uuid.NewString(),
		Source:      source,
		Type:        eventType,
		Time:        time.Now().UTC(),
		Data:        rawBytes,
	}, nil
}
