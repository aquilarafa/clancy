package ui

import (
	"strings"
	"testing"

	"github.com/aquila/clancy/model"
)

func TestRenderTextWrapsLongLines(t *testing.T) {
	// Text longer than width should wrap
	longText := "This is a very long line of text that should definitely wrap when rendered in a narrow terminal window"
	event := &model.DisplayEvent{
		Type: "assistant",
		Text: longText,
	}

	width := 40
	result := renderText(event, width)

	lines := strings.Split(result, "\n")
	for i, line := range lines {
		// Account for ANSI codes in length check - just verify no line is absurdly long
		if len(line) > width*3 { // generous margin for ANSI
			t.Errorf("line %d too long: %d chars", i, len(line))
		}
	}

	// Should have multiple lines due to wrapping
	if len(lines) < 2 {
		t.Errorf("expected wrapped lines, got %d line(s)", len(lines))
	}
}

func TestRenderToolUseWrapsInput(t *testing.T) {
	event := &model.DisplayEvent{
		Type: "assistant",
		ToolUse: &model.ToolUse{
			Name:  "Read",
			Input: `{"file_path": "/very/long/path/to/some/deeply/nested/directory/structure/file.go"}`,
		},
	}

	width := 50
	result := renderToolUse(event, width)

	if result == "" {
		t.Error("expected non-empty result")
	}

	// Should contain tool name
	if !strings.Contains(result, "Read") {
		t.Error("expected tool name in output")
	}
}

func TestRenderUserWrapsText(t *testing.T) {
	longText := "Please help me understand this complex codebase that has many interconnected components"
	event := &model.DisplayEvent{
		Type: "user",
		Text: longText,
	}

	width := 40
	result := renderUser(event, width)

	if result == "" {
		t.Error("expected non-empty result")
	}

	// Should start with prompt indicator
	if !strings.Contains(result, ">") {
		t.Error("expected > prompt in user output")
	}
}

func TestRenderThinkingWrapsText(t *testing.T) {
	event := &model.DisplayEvent{
		Type: "thinking",
		Text: "Let me think about this problem carefully and consider all the different approaches",
	}

	width := 40
	result := renderThinking(event, width)

	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestRenderToolResultWrapsContent(t *testing.T) {
	event := &model.DisplayEvent{
		Type: "tool_result",
		ToolResult: &model.ToolResult{
			ToolUseID: "123",
			Content:   "This is the content of a file that was read from disk and contains some useful information",
		},
	}

	width := 50
	result := renderToolResult(event, width)

	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestRenderEventDispatchesCorrectly(t *testing.T) {
	tests := []struct {
		name  string
		event *model.DisplayEvent
	}{
		{"system", &model.DisplayEvent{Type: "system", Text: "init"}},
		{"assistant", &model.DisplayEvent{Type: "assistant", Text: "hello"}},
		{"thinking", &model.DisplayEvent{Type: "thinking", Text: "hmm"}},
		{"user", &model.DisplayEvent{Type: "user", Text: "hi"}},
		{"result", &model.DisplayEvent{Type: "result", Text: "done"}},
		{"unknown", &model.DisplayEvent{Type: "unknown", Text: "wat"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderEvent(tt.event, 80)
			if result == "" && tt.event.Text != "" {
				t.Errorf("expected non-empty result for %s", tt.name)
			}
		})
	}
}

func TestRenderTodoWriteInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		empty    bool
	}{
		{
			name:     "all statuses",
			input:    `{"todos":[{"content":"Task pending","status":"pending"},{"content":"Task in progress","status":"in_progress"},{"content":"Task done","status":"completed"}]}`,
			contains: []string{"□", "Task pending", "◐", "Task in progress", "✓", "Task done"},
		},
		{
			name:     "pending only",
			input:    `{"todos":[{"content":"Do something","status":"pending"}]}`,
			contains: []string{"□", "Do something"},
		},
		{
			name:  "invalid json",
			input: `{invalid json`,
			empty: true,
		},
		{
			name:  "empty todos array",
			input: `{"todos":[]}`,
			empty: true,
		},
		{
			name:  "missing todos key",
			input: `{"other":"data"}`,
			empty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderTodoWriteInput(tt.input, 80)

			if tt.empty {
				if result != "" {
					t.Errorf("expected empty result, got: %s", result)
				}
				return
			}

			if result == "" {
				t.Error("expected non-empty result")
				return
			}

			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected result to contain %q", s)
				}
			}
		})
	}
}

func TestRenderToolUseTodoWrite(t *testing.T) {
	event := &model.DisplayEvent{
		Type: "assistant",
		ToolUse: &model.ToolUse{
			Name:  "TodoWrite",
			Input: `{"todos":[{"content":"First task","status":"pending"},{"content":"Second task","status":"in_progress"}]}`,
		},
	}

	result := renderToolUse(event, 80)

	if !strings.Contains(result, "TodoWrite") {
		t.Error("expected tool name in output")
	}
	if !strings.Contains(result, "□") {
		t.Error("expected pending icon □")
	}
	if !strings.Contains(result, "◐") {
		t.Error("expected in_progress icon ◐")
	}
	if !strings.Contains(result, "First task") {
		t.Error("expected first task content")
	}
	if !strings.Contains(result, "Second task") {
		t.Error("expected second task content")
	}
}

func TestRenderToolUseTodoWriteFallback(t *testing.T) {
	// Invalid JSON should fallback to normal rendering
	event := &model.DisplayEvent{
		Type: "assistant",
		ToolUse: &model.ToolUse{
			Name:  "TodoWrite",
			Input: `{invalid`,
		},
	}

	result := renderToolUse(event, 80)

	if !strings.Contains(result, "TodoWrite") {
		t.Error("expected tool name in output")
	}
	// Should show raw input as fallback
	if !strings.Contains(result, "{invalid") {
		t.Error("expected fallback to raw input")
	}
}
