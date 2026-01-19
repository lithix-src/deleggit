package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/datacraft/catalyst/core/internal/domain"
	"github.com/point-unknown/catalyst/pkg/logger"
	"github.com/point-unknown/catalyst/pkg/mcp"
	"github.com/point-unknown/catalyst/pkg/vector"
)

// LiaisonAgent is the Human-Swarm Interface.
// It translates user chat into structured commands (Git Issues).
type LiaisonAgent struct {
	id       string
	llm      domain.LLMProvider
	logger   *slog.Logger
	registry mcp.Registry
	vector   vector.Store
}

func NewLiaisonAgent(id string, llm domain.LLMProvider, reg mcp.Registry, vec vector.Store) *LiaisonAgent {
	return &LiaisonAgent{
		id:       id,
		llm:      llm,
		logger:   logger.New(fmt.Sprintf("agent-%s", id)),
		registry: reg,
		vector:   vec,
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

	// 0. Vibe Engine (RAG) retrieval
	// We allow the 'user_input' to drive the search.
	if userInput, ok := data["user_input"].(string); ok && userInput != "" && a.vector != nil {
		a.logger.Info("🤔 Recalling memories...", "query", userInput)

		// A. Embed
		embedding, err := a.llm.Embed(ctx, userInput)
		if err != nil {
			a.logger.Warn("Failed to embed user input", "error", err)
		} else {
			// B. Search
			results, err := a.vector.Search(ctx, embedding, 3)
			if err != nil {
				a.logger.Warn("Failed to search memories", "error", err)
			} else {
				a.logger.Info("found memories", "count", len(results))
				for _, r := range results {
					a.logger.Info("Memory", "score", r.Score, "preview", r.Content[:min(len(r.Content), 50)]+"...")
				}
				// TODO: Append results to LLM Context
			}
		}
	}

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
