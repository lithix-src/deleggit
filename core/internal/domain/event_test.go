package domain

import (
	"encoding/json"
	"testing"
)

func TestNewEvent(t *testing.T) {
	source := "unit-test"
	eventType := "test.event"
	payload := map[string]string{"foo": "bar"}

	evt, err := NewEvent(source, eventType, payload)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	if evt.Source != source {
		t.Errorf("Expected Source %s, got %s", source, evt.Source)
	}
	if evt.Type != eventType {
		t.Errorf("Expected Type %s, got %s", eventType, evt.Type)
	}
	if evt.SpecVersion != "1.0" {
		t.Errorf("Expected SpecVersion 1.0, got %s", evt.SpecVersion)
	}
}

func TestCloudEvent_Marshaling(t *testing.T) {
	// 1. Create Event
	payload := map[string]int{"value": 42}
	evt, _ := NewEvent("src", "type", payload)

	// 2. Marshal to JSON
	bytes, err := json.Marshal(evt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 3. Unmarshal back
	var loaded evtCopy
	if err := json.Unmarshal(bytes, &loaded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if loaded.Source != "src" {
		t.Errorf("Expected Source 'src', got %s", loaded.Source)
	}
}

// Helper struct to avoid method conflicts in test
type evtCopy struct {
	ID     string          `json:"id"`
	Source string          `json:"source"`
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
}
