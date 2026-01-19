package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/datacraft/catalyst/core/internal/domain"
	"github.com/point-unknown/catalyst/pkg/logger"
	"github.com/point-unknown/catalyst/pkg/mcp"
)

// LiaisonAgent is the Human-Swarm Interface.
// It translates user chat into structured commands (Git Issues).
type LiaisonAgent struct {
	id       string
	llm      domain.LLMProvider
	logger   *slog.Logger
	registry mcp.Registry
}

func NewLiaisonAgent(id string, llm domain.LLMProvider, reg mcp.Registry) *LiaisonAgent {
	return &LiaisonAgent{
		id:       id,
		llm:      llm,
		logger:   logger.New(fmt.Sprintf("agent-%s", id)),
		registry: reg,
	}
}

func (a *LiaisonAgent) ID() string {
	return a.id
}

func (a *LiaisonAgent) Type() domain.AgentType {
	return domain.AgentTypeCommunicator // Level 2: Communicator
}

func (a *LiaisonAgent) Execute(ctx context.Context, input domain.CloudEvent) (*domain.CloudEvent, error) {
	// Trigger: "agent.{id}.wake" from GroupChatManager
	if input.Type != fmt.Sprintf("agent.%s.wake", a.id) {
		return nil, nil
	}

	// Unmarshal context data
	var data map[string]interface{}
	if len(input.Data) > 0 {
		if err := json.Unmarshal(input.Data, &data); err != nil {
			a.logger.Warn("Failed to unmarshal wake data", "error", err)
			data = make(map[string]interface{})
		}
	} else {
		data = make(map[string]interface{})
	}

	a.logger.Info("Liaison woke up. Listening to chat context...", "context_len", data["context_len"])

	// 1. List Available Tools
	tools := a.registry.ListTools()
	toolsListJSON, _ := json.Marshal(tools)
	// Log the tools for debugging, but use the underscore to avoid "unused variable" error if we don't use it elsewhere yet.
	_ = toolsListJSON
	a.logger.Info("Available Tools", "count", len(tools))

	// 2. "Think" (Simulate LLM tool selection)
	// In reality, we pass 'toolsListJSON' to the LLM system prompt.
	// For MVP, if input context implies "Refactor", we pick 'git_create_issue'.

	// Simulated Decision:
	triggerTool := true // Logic would go here

	if triggerTool {
		toolName := "git_create_issue"
		// Check if tool exists
		if _, ok := a.registry.GetTool(toolName); ok {
			a.logger.Info("Decided to call tool", "tool", toolName)

			// Emit Tool Call
			toolCall := mcp.ToolCall{
				ID:       "call_12345",
				ToolName: toolName,
				Arguments: map[string]interface{}{
					"title": "User Request: Refactor Login",
					"body":  "As requested in ChatOps session.",
				},
			}

			// We return the Tool Call event.
			// Type: "tool.call"
			evt, _ := domain.NewEvent(a.id, "tool.call", toolCall)

			return &evt, nil
		}
	}

	responseMsg := "I'm listening, but I didn't see a need to open a ticket."
	reply, _ := domain.NewEvent(a.id, "chat.message", map[string]string{
		"content": responseMsg,
		"sender":  "Liaison",
	})

	return &reply, nil
}
