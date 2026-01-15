package model

import "encoding/json"

// Event represents a Claude Code stream-json event
type Event struct {
	Type      string          `json:"type"` // system, assistant, user, result
	Subtype   string          `json:"subtype,omitempty"`
	Message   *Message        `json:"message,omitempty"`
	SessionID string          `json:"session_id,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
	Cwd       string          `json:"cwd,omitempty"`
	GitBranch string          `json:"gitBranch,omitempty"`
	Raw       json.RawMessage `json:"-"`

	// Result fields
	CostUSD     float64 `json:"cost_usd,omitempty"`
	DurationMS  int     `json:"duration_ms,omitempty"`
	NumTurns    int     `json:"num_turns,omitempty"`

	// System init fields
	Tools []string `json:"tools,omitempty"`
	Model string   `json:"model,omitempty"`
}

// Message is the assistant/user message structure
type Message struct {
	ID         string          `json:"id,omitempty"`
	Type       string          `json:"type,omitempty"`
	Role       string          `json:"role"`
	Model      string          `json:"model,omitempty"`
	Content    json.RawMessage `json:"content"` // can be string or array
	StopReason *string         `json:"stop_reason,omitempty"`
	Usage      *Usage          `json:"usage,omitempty"`
}

// ContentBlock is an item in message.content array
type ContentBlock struct {
	Type      string          `json:"type"` // text, tool_use, tool_result, thinking
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`   // for tool_use
	Name      string          `json:"name,omitempty"` // for tool_use
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"` // for tool_result
	Content   json.RawMessage `json:"content,omitempty"`     // for tool_result
	Thinking  string          `json:"thinking,omitempty"`
}

// Usage contains token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// DisplayEvent is a processed event ready for rendering
type DisplayEvent struct {
	Type       string // system, assistant, user, thinking, tool_result, result
	Text       string
	ToolUse    *ToolUse
	ToolResult *ToolResult
	Model      string
	Cwd        string
	Usage      *Usage
	StopReason string
	CostUSD    float64
	NumTurns   int
	DurationMS int
}

// ToolUse represents a tool invocation
type ToolUse struct {
	ID    string
	Name  string
	Input string // JSON string of input
}

// ToolResult represents a tool result
type ToolResult struct {
	ToolUseID string
	Content   string // truncated content
}
