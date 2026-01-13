package service

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/datacraft/catalyst/core/internal/domain"
)

// Mock Registry Agent
type RegistryAgent struct {
	id string
}

func (r *RegistryAgent) ID() string             { return r.id }
func (r *RegistryAgent) Type() domain.AgentType { return domain.AgentTypeReporter }
func (r *RegistryAgent) Execute(ctx context.Context, e domain.CloudEvent) (*domain.CloudEvent, error) {
	return nil, nil
}

func TestAgentRegistry_Lifecycle(t *testing.T) {
	reg := NewAgentRegistry()
	agent := &RegistryAgent{id: "TestAgent"}

	// 1. Register
	reg.Register(agent)

	// 2. Get
	got, err := reg.Get("TestAgent")
	if err != nil {
		t.Fatalf("Failed to get agent: %v", err)
	}
	if got.ID() != "TestAgent" {
		t.Errorf("Expected ID 'TestAgent', got %s", got.ID())
	}

	// 3. Get Non-Existent
	_, err = reg.Get("Ghost")
	if err == nil {
		t.Errorf("Expected error for missing agent, got nil")
	}

	// 4. List
	list := reg.List()
	if len(list) != 1 {
		t.Errorf("Expected 1 agent in list, got %d", len(list))
	}
}

func TestAgentRegistry_Concurrency(t *testing.T) {
	reg := NewAgentRegistry()
	var wg sync.WaitGroup

	// Concurrent Writers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("Agent-%d", n)
			reg.Register(&RegistryAgent{id: id})
		}(i)
	}

	// Concurrent Readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reg.List()
		}()
	}

	wg.Wait()

	// Verify count
	list := reg.List()
	if len(list) != 100 {
		t.Errorf("Expected 100 agents after concurrent writes, got %d", len(list))
	}
}
