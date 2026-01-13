package agent

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/datacraft/catalyst/core/internal/domain"
)

func TestTrendScout_Execute(t *testing.T) {
	// Setup
	agent := NewTrendScout("TestScout", 80.0)
	ctx := context.Background()

	// Scenario 1: Ignore non-matching event
	evt1, _ := domain.NewEvent("sensor-1", "other.type", nil)
	out1, err := agent.Execute(ctx, evt1)
	if err != nil {
		t.Errorf("Unexpected error on ignore: %v", err)
	}
	if out1 != nil {
		t.Errorf("Expected nil output for non-matching event")
	}

	// Scenario 2: Normal Data (Below Threshold)
	// Send 5 readings of 50.0. Avg = 50.0 < 80.0
	for i := 0; i < 5; i++ {
		payload := map[string]interface{}{"value": 50.0, "unit": "C"}
		data, _ := json.Marshal(payload)
		evt := domain.CloudEvent{Type: "sensor.cpu.temp", Data: data, Time: time.Now()}

		out, err := agent.Execute(ctx, evt)
		if err != nil {
			t.Fatalf("Execution failed: %v", err)
		}
		if out != nil {
			t.Errorf("Expected no alert for temp 50.0")
		}
	}

	// Scenario 3: Spike (Overheat)
	// Send reading of 1000.0. Window is 5.
	// Previous: [50, 50, 50, 50, 50] -> Avg 50
	// New: [50, 50, 50, 50, 1000] -> Avg 240 > 80
	payload := map[string]interface{}{"value": 1000.0, "unit": "C"}
	data, _ := json.Marshal(payload)
	evt := domain.CloudEvent{Type: "sensor.cpu.temp", Data: data, Time: time.Now()}

	out, err := agent.Execute(ctx, evt)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}
	if out == nil {
		t.Fatalf("Expected Alert, got nil")
	}

	if out.Type != "swarm.security.alert" {
		t.Errorf("Expected signal type 'swarm.security.alert', got %s", out.Type)
	}
}
