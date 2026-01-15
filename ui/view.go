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
		return renderSystem(event)
	case "assistant":
		if event.ToolUse != nil {
			return renderToolUse(event)
		}
		return renderText(event)
	case "thinking":
		return renderThinking(event)
	case "tool_result":
		return renderToolResult(event)
	case "user":
		return renderUser(event)
	case "result":
		return renderResult(event)
	default:
		return renderUnknown(event)
	}
}

func renderSystem(event *model.DisplayEvent) string {
	badge := badgeSystem.Render("system")
	if event.Text != "" {
		return eventStyle.Render(fmt.Sprintf("%s\n  %s", badge, usageStyle.Render(event.Text)))
	}
	return eventStyle.Render(badge)
}

func renderText(event *model.DisplayEvent) string {
	badge := badgeText.Render("text")
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
	text = strings.Join(lines, "\n  ")
	return eventStyle.Render(fmt.Sprintf("%s\n  %s", badge, textStyle.Render(text)))
}

func renderThinking(event *model.DisplayEvent) string {
	badge := badgeThinking.Render("thinking")
	text := event.Text
	if len(text) > 200 {
		text = text[:200] + "..."
	}
	text = strings.TrimSpace(text)
	return eventStyle.Render(fmt.Sprintf("%s\n  %s", badge, thinkingStyle.Render(text)))
}

func renderToolUse(event *model.DisplayEvent) string {
	badge := badgeTool.Render("tool_use")
	tool := event.ToolUse
	if tool == nil {
		return eventStyle.Render(badge)
	}

	toolName := toolNameStyle.Render(tool.Name)
	input := tool.Input
	if len(input) > 150 {
		input = input[:150] + "..."
	}
	return eventStyle.Render(fmt.Sprintf("%s %s\n  %s", badge, toolName, toolInputStyle.Render(input)))
}

func renderUser(event *model.DisplayEvent) string {
	badge := badgeUser.Render("user")
	text := event.Text
	if len(text) > 200 {
		text = text[:200] + "..."
	}
	text = strings.TrimSpace(text)
	return eventStyle.Render(fmt.Sprintf("%s\n  %s", badge, textStyle.Render(text)))
}

func renderToolResult(event *model.DisplayEvent) string {
	badge := badgeResult.Render("result")
	if event.ToolResult == nil {
		return eventStyle.Render(badge)
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
	return eventStyle.Render(fmt.Sprintf("%s\n  %s", badge, resultStyle.Render(content)))
}

func renderResult(event *model.DisplayEvent) string {
	badge := badgeSuccess.Render("done")
	return eventStyle.Render(fmt.Sprintf("%s %s", badge, successStyle.Render(event.Text)))
}

func renderUnknown(event *model.DisplayEvent) string {
	badge := badgeDefault.Render(event.Type)
	if event.Text != "" {
		text := event.Text
		if len(text) > 100 {
			text = text[:100] + "..."
		}
		return eventStyle.Render(fmt.Sprintf("%s\n  %s", badge, text))
	}
	return eventStyle.Render(badge)
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
