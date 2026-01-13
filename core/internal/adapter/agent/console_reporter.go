package agent

import (
	"context"
	"log"

	"github.com/datacraft/deleggit/core/internal/domain"
)

// ConsoleReporter is a Phase 2 test agent that simply reports activity to stdout.
// It verifies the "Reporting" implementation of the spectrum.
type ConsoleReporter struct {
	id string
}

func NewConsoleReporter(id string) *ConsoleReporter {
	return &ConsoleReporter{id: id}
}

func (c *ConsoleReporter) ID() string {
	return c.id
}

func (c *ConsoleReporter) Type() domain.AgentType {
	return domain.AgentTypeReporter
}

func (c *ConsoleReporter) Execute(ctx context.Context, input domain.CloudEvent) (*domain.CloudEvent, error) {
	log.Printf("  >>> [AGENT:%s] Processing Event: %s", c.id, input.Type)
	// No output for a reporter
	return nil, nil
}
