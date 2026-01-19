package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/datacraft/catalyst/core/internal/domain"
)

// EngineerAgent is a "Doer" (Level 3) capable of fixing code.
type EngineerAgent struct {
	id        string
	llm       domain.LLMProvider
	workspace domain.Workspace
	security  domain.SecurityConfig
}

func NewEngineerAgent(id string, llm domain.LLMProvider, ws domain.Workspace, security domain.SecurityConfig) *EngineerAgent {
	return &EngineerAgent{
		id:        id,
		llm:       llm,
		workspace: ws,
		security:  security,
	}
}

func (a *EngineerAgent) ID() string {
	return a.id
}

func (a *EngineerAgent) Type() domain.AgentType {
	return domain.AgentTypeExpressor // Level 3: Can produce artifacts (Create PRs)
}

func (a *EngineerAgent) Execute(ctx context.Context, input domain.CloudEvent) (*domain.CloudEvent, error) {
	// Trigger: "repo.issue.command"
	if input.Type != "repo.issue.command" {
		// The original code returned nil, nil.
		// To make `return &evt, nil` syntactically correct without
		// introducing unrelated edits, we must define `evt`.
		// Assuming `evt` is intended to be an empty or default CloudEvent
		// when the input type doesn't match.
		// If the intent was to return the input event, it would be `&input`.
		// If the intent was to return nothing, `nil, nil` is appropriate.
		// For now, we define a zero-value CloudEvent to satisfy the syntax.
		evt := domain.CloudEvent{}
		return &evt, nil
	}

	log.Printf("[ENGINEER:%s] analyzing issue...", a.id)

	// In Phase 5.2 We just "Think" (Generate Plan)
	// Future: Write Code

	// Mocking extraction of context
	// In reality this comes from 'input.Data'
	issueContext := fmt.Sprintf("Analyze input: %v", input.Data)

	plan, err := a.llm.GenerateCode(ctx, "Create a fix plan for: "+issueContext)
	if err != nil {
		log.Printf("[ENGINEER:%s] Brain Freeze: %v", a.id, err)
		return nil, err
	}

	log.Printf("[ENGINEER:%s] Generated Plan: %s", a.id, plan)

	// Return the Plan as an Event
	evt, err := domain.NewEvent("agent.engineer", "agent.plan.generated", map[string]string{
		"agent": a.id,
		"plan":  plan,
	})
	if err != nil {
		return nil, err
	}
	return &evt, nil
}
