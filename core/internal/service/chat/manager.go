package chat

import (
	"fmt"
	"log/slog"

	"github.com/datacraft/catalyst/core/internal/domain"
)

// GroupChatManager orchestrates the conversation.
// Based on AutoGen's GroupChat pattern.
type GroupChatManager struct {
	history []domain.CloudEvent
	logger  *slog.Logger
	pub     func(topic string, event domain.CloudEvent)
}

func NewGroupChatManager(logger *slog.Logger, publisher func(topic string, event domain.CloudEvent)) *GroupChatManager {
	return &GroupChatManager{
		history: make([]domain.CloudEvent, 0),
		logger:  logger,
		pub:     publisher,
	}
}

// HandleUserMessage processes input from the Mission Control UI.
func (m *GroupChatManager) HandleUserMessage(msg string) {
	// 1. Record User Message
	evt, _ := domain.NewEvent("user", "chat.message", map[string]string{"content": msg})
	m.addToHistory(evt)

	// 2. Select Next Speaker (Round Robin: User -> Liaison)
	// In the future, this will be smarter (Auto).
	m.selectNextSpeaker("liaison")
}

func (m *GroupChatManager) addToHistory(evt domain.CloudEvent) {
	m.history = append(m.history, evt)
	// Broadcast to UI
	m.pub("chat/stream/message", evt)
}

func (m *GroupChatManager) selectNextSpeaker(agentID string) {
	// Trigger the agent to speak
	// We use specific type "agent.liaison.wake" so MissionManager can route it exactly.
	eventType := fmt.Sprintf("agent.%s.wake", agentID)
	cmd, _ := domain.NewEvent("chat-manager", eventType, map[string]string{
		"target_agent": agentID,
		"context_len":  fmt.Sprintf("%d", len(m.history)),
	})
	m.pub(fmt.Sprintf("agent/%s/wake", agentID), cmd)
}

// ProcessAgentResponse handles when an agent speaks back.
func (m *GroupChatManager) ProcessAgentResponse(evt domain.CloudEvent) {
	m.addToHistory(evt)
	// Logic to decide if we stop or pick another speaker would go here.
}
