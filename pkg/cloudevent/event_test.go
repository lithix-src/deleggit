package cloudevent

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	evt, err := New("test-source", "test.type", data)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if evt.Source != "test-source" {
		t.Errorf("New() Source = %v, want %v", evt.Source, "test-source")
	}
	if evt.Type != "test.type" {
		t.Errorf("New() Type = %v, want %v", evt.Type, "test.type")
	}
	if evt.SpecVersion != "1.0" {
		t.Errorf("New() SpecVersion = %v, want %v", evt.SpecVersion, "1.0")
	}
	if evt.ID == "" {
		t.Error("New() ID is empty")
	}
	if evt.Time.IsZero() {
		t.Error("New() Time is zero")
	}

	// Verify Data Marshalling
	var parsedData map[string]string
	if err := json.Unmarshal(evt.Data, &parsedData); err != nil {
		t.Fatalf("Failed to unmarshal data: %v", err)
	}
	if parsedData["foo"] != "bar" {
		t.Errorf("Data['foo'] = %v, want %v", parsedData["foo"], "bar")
	}
}
