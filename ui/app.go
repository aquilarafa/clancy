package ui

import (
	"fmt"
	"strings"

	"github.com/aquila/clancy/model"
	"github.com/aquila/clancy/parser"
	"github.com/aquila/clancy/watcher"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the main bubbletea model
type Model struct {
	filename   string
	watcher    *watcher.Watcher
	parser     *parser.Parser
	events     []*model.DisplayEvent
	width      int
	height     int
	offset     int // scroll offset
	followMode bool
	err        error
}

// lineMsg is a message containing a new line from the watcher
type lineMsg []byte

// errMsg is a message containing an error
type errMsg error

// New creates a new UI model
func New(filename string, w *watcher.Watcher) Model {
	return Model{
		filename:   filename,
		watcher:    w,
		parser:     parser.New(),
		events:     make([]*model.DisplayEvent, 0),
		followMode: true,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		waitForLine(m.watcher),
		waitForError(m.watcher),
	)
}

// waitForLine waits for the next line from the watcher
func waitForLine(w *watcher.Watcher) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-w.Lines()
		if !ok {
			return nil
		}
		return lineMsg(line)
	}
}

// waitForError waits for errors from the watcher
func waitForError(w *watcher.Watcher) tea.Cmd {
	return func() tea.Msg {
		err, ok := <-w.Errors()
		if !ok {
			return nil
		}
		return errMsg(err)
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.watcher.Stop()
			return m, tea.Quit

		case "up", "k":
			if m.offset > 0 {
				m.offset--
				m.followMode = false
			}

		case "down", "j":
			maxOffset := m.maxOffset()
			if m.offset < maxOffset {
				m.offset++
			}

		case "g", "home":
			m.offset = 0
			m.followMode = false

		case "G", "end":
			m.offset = m.maxOffset()
			m.followMode = true

		case "f":
			m.followMode = !m.followMode
			if m.followMode {
				m.offset = m.maxOffset()
			}

		case "pgup":
			m.offset -= m.viewportHeight()
			if m.offset < 0 {
				m.offset = 0
			}
			m.followMode = false

		case "pgdown":
			m.offset += m.viewportHeight()
			maxOffset := m.maxOffset()
			if m.offset > maxOffset {
				m.offset = maxOffset
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case lineMsg:
		events, err := m.parser.ParseLine(msg)
		if err == nil && len(events) > 0 {
			m.events = append(m.events, events...)
			if m.followMode {
				m.offset = m.maxOffset()
			}
		}
		return m, waitForLine(m.watcher)

	case errMsg:
		m.err = msg
		return m, waitForError(m.watcher)
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Status bar
	b.WriteString(renderStatusBar(m.filename, len(m.events), m.width))
	b.WriteString("\n")

	// Viewport content
	viewportHeight := m.viewportHeight()
	content := m.renderEvents()
	lines := strings.Split(content, "\n")

	// Apply scroll offset
	start := m.offset
	end := start + viewportHeight
	if end > len(lines) {
		end = len(lines)
	}
	if start > len(lines) {
		start = len(lines)
	}

	visibleLines := lines[start:end]

	// Pad to fill viewport
	for len(visibleLines) < viewportHeight {
		visibleLines = append(visibleLines, "")
	}

	b.WriteString(strings.Join(visibleLines, "\n"))
	b.WriteString("\n")

	// Help bar
	b.WriteString(renderHelpBar(m.followMode, m.width))

	return b.String()
}

// renderEvents renders all events
func (m Model) renderEvents() string {
	var parts []string
	for _, event := range m.events {
		rendered := renderEvent(event, m.width)
		if rendered != "" {
			parts = append(parts, rendered)
		}
	}
	if len(parts) == 0 {
		return fmt.Sprintf("\n  Waiting for events from %s...\n", m.filename)
	}
	return strings.Join(parts, "\n")
}

// viewportHeight returns the height available for events
func (m Model) viewportHeight() int {
	// Total height minus status bar (1) and help bar (1)
	h := m.height - 2
	if h < 1 {
		h = 1
	}
	return h
}

// maxOffset returns the maximum scroll offset
func (m Model) maxOffset() int {
	content := m.renderEvents()
	lines := strings.Split(content, "\n")
	max := len(lines) - m.viewportHeight()
	if max < 0 {
		return 0
	}
	return max
}
