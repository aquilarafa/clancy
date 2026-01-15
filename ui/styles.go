package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Adaptive colors following charmbracelet/bubbles conventions
	// Light/Dark pairs for terminal background adaptation
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	muted  = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#626262"}
	text   = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}
	accent = lipgloss.AdaptiveColor{Light: "#7571F9", Dark: "#7571F9"}
	err    = lipgloss.AdaptiveColor{Light: "#FF5F87", Dark: "#FF5F87"}

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Background(subtle).
			Foreground(text).
			Padding(0, 1)

	// Help bar
	helpBarStyle = lipgloss.NewStyle().
			Foreground(muted).
			Padding(0, 1)

	// Badge base style
	badgeStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)

	// Event type badges
	badgeSystem = badgeStyle.
			Background(muted).
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"})

	badgeUser = badgeStyle.
			Background(text).
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"})

	badgeText = badgeStyle.
			Background(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#000000"})

	badgeTool = badgeStyle.
			Background(accent).
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"})

	badgeThinking = badgeStyle.
			Background(subtle).
			Foreground(muted)

	badgeResult = badgeStyle.
			Background(muted).
			Foreground(text)

	badgeSuccess = badgeStyle.
			Background(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"})

	badgeError = badgeStyle.
			Background(err).
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"})

	badgeDefault = badgeStyle.
			Background(subtle).
			Foreground(text)

	// Content styles
	textStyle = lipgloss.NewStyle().
			Foreground(text)

	toolNameStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	toolInputStyle = lipgloss.NewStyle().
			Foreground(muted)

	thinkingStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(text)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})

	errorStyle = lipgloss.NewStyle().
			Foreground(err)

	usageStyle = lipgloss.NewStyle().
			Foreground(muted)

	// Event container
	eventStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			MarginBottom(1)

	// Follow mode indicator
	followOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
			Bold(true)

	followOffStyle = lipgloss.NewStyle().
			Foreground(muted)
)
