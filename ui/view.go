package ui

import (
	"fmt"
	"strings"

	"github.com/aquila/clancy/model"
)

// renderEvent renders a single display event
func renderEvent(event *model.DisplayEvent, width int) string {
	switch event.Type {
	case "system":
		return renderSystem(event, width)
	case "assistant":
		if event.ToolUse != nil {
			return renderToolUse(event, width)
		}
		return renderText(event, width)
	case "thinking":
		return renderThinking(event, width)
	case "tool_result":
		return renderToolResult(event, width)
	case "user":
		return renderUser(event, width)
	case "result":
		return renderResult(event, width)
	default:
		return renderUnknown(event, width)
	}
}

func renderSystem(event *model.DisplayEvent, width int) string {
	if event.Text != "" {
		contentWidth := width - 4 // account for padding
		return eventStyle.Width(width).Render(usageStyle.Width(contentWidth).Render(event.Text))
	}
	return ""
}

func renderText(event *model.DisplayEvent, width int) string {
	text := event.Text
	if len(text) > 300 {
		text = text[:300] + "..."
	}
	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")
	if len(lines) > 5 {
		lines = lines[:5]
		lines = append(lines, "...")
	}
	text = strings.Join(lines, "\n")
	contentWidth := width - 4
	return eventStyle.Width(width).Render(textStyle.Width(contentWidth).Render(text))
}

func renderThinking(event *model.DisplayEvent, width int) string {
	text := event.Text
	if len(text) > 200 {
		text = text[:200] + "..."
	}
	text = strings.TrimSpace(text)
	contentWidth := width - 4
	return eventStyle.Width(width).Render(thinkingStyle.Width(contentWidth).Render(text))
}

func renderToolUse(event *model.DisplayEvent, width int) string {
	tool := event.ToolUse
	if tool == nil {
		return ""
	}

	toolName := toolNameStyle.Render("● " + tool.Name)
	input := tool.Input
	if len(input) > 150 {
		input = input[:150] + "..."
	}
	contentWidth := width - 6
	return eventStyle.Width(width).Render(fmt.Sprintf("%s\n  %s", toolName, toolInputStyle.Width(contentWidth).Render(input)))
}

func renderUser(event *model.DisplayEvent, width int) string {
	text := event.Text
	if len(text) > 200 {
		text = text[:200] + "..."
	}
	text = strings.TrimSpace(text)
	contentWidth := width - 6
	return eventStyle.Width(width).Render(fmt.Sprintf("> %s", textStyle.Width(contentWidth).Render(text)))
}

func renderToolResult(event *model.DisplayEvent, width int) string {
	if event.ToolResult == nil {
		return ""
	}

	content := event.ToolResult.Content
	if len(content) > 200 {
		content = content[:200] + "..."
	}
	lines := strings.Split(content, "\n")
	if len(lines) > 4 {
		lines = lines[:4]
		lines = append(lines, "...")
	}
	content = strings.Join(lines, "\n  ")
	contentWidth := width - 6
	return eventStyle.Width(width).Render(fmt.Sprintf("  %s", resultStyle.Width(contentWidth).Render(content)))
}

func renderResult(event *model.DisplayEvent, width int) string {
	contentWidth := width - 4
	return eventStyle.Width(width).Render(successStyle.Width(contentWidth).Render("✓ " + event.Text))
}

func renderUnknown(event *model.DisplayEvent, width int) string {
	if event.Text != "" {
		text := event.Text
		if len(text) > 100 {
			text = text[:100] + "..."
		}
		contentWidth := width - 4
		return eventStyle.Width(width).Render(textStyle.Width(contentWidth).Render(text))
	}
	return ""
}

// renderStatusBar renders the top status bar
func renderStatusBar(filename string, eventCount int, width int) string {
	left := fmt.Sprintf(" watching: %s", filename)
	right := fmt.Sprintf("%d events ", eventCount)
	spaces := width - len(left) - len(right)
	if spaces < 1 {
		spaces = 1
	}
	return statusBarStyle.Width(width).Render(left + strings.Repeat(" ", spaces) + right)
}

// renderHelpBar renders the bottom help bar
func renderHelpBar(followMode bool, width int) string {
	followIndicator := ""
	if followMode {
		followIndicator = followOnStyle.Render("[FOLLOW]")
	} else {
		followIndicator = followOffStyle.Render("[follow off]")
	}
	help := fmt.Sprintf("q:quit  ↑↓/jk:scroll  g/G:top/bottom  f:follow  %s", followIndicator)
	return helpBarStyle.Width(width).Render(help)
}
