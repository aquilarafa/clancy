package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aquila/clancy/model"
)

// Parser processes NDJSON lines from Claude Code stream-json output
type Parser struct{}

// New creates a new Parser
func New() *Parser {
	return &Parser{}
}

// ParseLine parses a single JSON line and returns DisplayEvents
func (p *Parser) ParseLine(line []byte) ([]*model.DisplayEvent, error) {
	if len(line) == 0 {
		return nil, nil
	}

	var event model.Event
	if err := json.Unmarshal(line, &event); err != nil {
		return nil, err
	}
	event.Raw = line

	var events []*model.DisplayEvent

	switch event.Type {
	case "system":
		events = append(events, &model.DisplayEvent{
			Type:  "system",
			Model: event.Model,
			Cwd:   event.Cwd,
			Text:  fmt.Sprintf("%s | %s", event.Model, event.Cwd),
		})

	case "assistant":
		if event.Message == nil {
			return nil, nil
		}
		blocks := p.parseContent(event.Message.Content)
		for _, block := range blocks {
			de := &model.DisplayEvent{
				Type:  "assistant",
				Model: event.Message.Model,
				Usage: event.Message.Usage,
			}
			if event.Message.StopReason != nil {
				de.StopReason = *event.Message.StopReason
			}

			switch block.Type {
			case "text":
				if block.Text != "" {
					de.Text = block.Text
					events = append(events, de)
				}
			case "tool_use":
				inputStr := string(block.Input)
				de.ToolUse = &model.ToolUse{
					ID:    block.ID,
					Name:  block.Name,
					Input: inputStr,
				}
				events = append(events, de)
			case "thinking":
				if block.Thinking != "" {
					de.Text = block.Thinking
					de.Type = "thinking"
					events = append(events, de)
				}
			}
		}

	case "user":
		if event.Message == nil {
			return nil, nil
		}
		// Check if content is a string (human input) or array (tool results)
		content := event.Message.Content
		if len(content) > 0 && content[0] == '"' {
			// String content - human input
			var text string
			if err := json.Unmarshal(content, &text); err == nil && text != "" {
				events = append(events, &model.DisplayEvent{
					Type: "user",
					Text: text,
				})
			}
		} else {
			// Array content - tool results
			blocks := p.parseContent(content)
			for _, block := range blocks {
				if block.Type == "tool_result" {
					contentStr := p.extractToolResultContent(block.Content)
					if len(contentStr) > 500 {
						contentStr = contentStr[:500] + "..."
					}
					events = append(events, &model.DisplayEvent{
						Type: "tool_result",
						ToolResult: &model.ToolResult{
							ToolUseID: block.ToolUseID,
							Content:   contentStr,
						},
					})
				}
			}
		}

	case "result":
		if event.Subtype == "success" {
			events = append(events, &model.DisplayEvent{
				Type:       "result",
				Text:       fmt.Sprintf("✓ %d turns | $%.4f | %dms", event.NumTurns, event.CostUSD, event.DurationMS),
				CostUSD:    event.CostUSD,
				NumTurns:   event.NumTurns,
				DurationMS: event.DurationMS,
			})
		}

	default:
		// Unknown event type - show type name
		if event.Type != "" {
			events = append(events, &model.DisplayEvent{
				Type: event.Type,
			})
		}
	}

	return events, nil
}

// parseContent handles content that can be string or array
func (p *Parser) parseContent(raw json.RawMessage) []model.ContentBlock {
	if len(raw) == 0 {
		return nil
	}

	// Try array first
	var blocks []model.ContentBlock
	if err := json.Unmarshal(raw, &blocks); err == nil {
		return blocks
	}

	// Try string
	var text string
	if err := json.Unmarshal(raw, &text); err == nil && text != "" {
		return []model.ContentBlock{{Type: "text", Text: text}}
	}

	return nil
}

// extractToolResultContent handles tool_result content that can be string or object
func (p *Parser) extractToolResultContent(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	// Try string first
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}

	// For objects, return compact JSON
	var obj interface{}
	if err := json.Unmarshal(raw, &obj); err == nil {
		b, _ := json.Marshal(obj)
		s := string(b)
		// Truncate long JSON
		if len(s) > 300 {
			s = s[:300] + "..."
		}
		return s
	}

	return string(raw)
}

// truncateLines truncates content to max lines
func truncateLines(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	return strings.Join(lines[:maxLines], "\n") + "\n..."
}
