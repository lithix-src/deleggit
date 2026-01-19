package mcp

import (
	"encoding/json"
)

// Tool represents a capability capability exposed to Agents.
// It follows a simplified MCP/OpenAI-function schema.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"` // JSON Schema defining inputs
}

// ToolCall represents a request from an Agent to execute a Tool.
type ToolCall struct {
	ID        string                 `json:"id"`
	ToolName  string                 `json:"tool_name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the output of a Tool execution.
type ToolResult struct {
	CallID string `json:"call_id"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// Registry is a local interface for managing tools.
type Registry interface {
	ListTools() []Tool
	GetTool(name string) (*Tool, bool)
}
