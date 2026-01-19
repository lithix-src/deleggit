package service

import (
	"context"
	"testing"

	"github.com/datacraft/catalyst/core/internal/domain"
)

// MockAgent implements domain.Agent for testing
type MockAgent struct {
	id          string
	executed    bool
	lastInput   domain.CloudEvent
	outputEvent *domain.CloudEvent
}

func (m *MockAgent) ID() string {
	return m.id
}
func (m *MockAgent) Type() domain.AgentType {
	return domain.AgentTypeReporter
}
func (m *MockAgent) Execute(ctx context.Context, input domain.CloudEvent) (*domain.CloudEvent, error) {
	m.executed = true
	m.lastInput = input
	return m.outputEvent, nil
}

func TestMissionManager_ProcessEvent(t *testing.T) {
	// Setup Registry
	registry := NewAgentRegistry()
	mockAgent := &MockAgent{id: "TestAgent"}
	registry.Register(mockAgent)

	// Setup Manager
	// Mock Publisher
	publisher := func(topic string, event domain.CloudEvent) {
		// Do nothing or record calls
	}
	manager := NewMissionManager(registry, publisher)

	// Load Mission: "sensor/test" -> "TestAgent"
	manager.LoadMission(domain.Mission{
		ID:           "m1",
		Name:         "Test Mission",
		TriggerTopic: "sensor/test",
		Agents:       []string{"TestAgent"},
	})

	// Scenario 1: Matching Event
	evt1, _ := domain.NewEvent("source1", "sensor/test", nil)
	manager.ProcessEvent(evt1)

	if !mockAgent.executed {
		t.Errorf("Expected agent to execute for matching triggering topic")
	}
	if mockAgent.lastInput.Type != "sensor/test" {
		t.Errorf("Expected agent to receive the event")
	}

	// Reset
	mockAgent.executed = false

	// Scenario 2: Non-matching Event
	evt2, _ := domain.NewEvent("source1", "sensor/other", nil)
	manager.ProcessEvent(evt2)

	if mockAgent.executed {
		t.Errorf("Agent executed for non-matching topic")
	}
}
