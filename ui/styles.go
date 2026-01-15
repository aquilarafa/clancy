package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorBlue    = lipgloss.Color("#5C9DFF")
	colorGreen   = lipgloss.Color("#5CFF9D")
	colorYellow  = lipgloss.Color("#FFD75C")
	colorRed     = lipgloss.Color("#FF5C5C")
	colorDim     = lipgloss.Color("#666666")
	colorWhite   = lipgloss.Color("#FFFFFF")
	colorCyan    = lipgloss.Color("#5CFFFF")
	colorMagenta = lipgloss.Color("#FF5CFF")
	colorOrange  = lipgloss.Color("#FFA55C")

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(colorWhite).
			Padding(0, 1)

	// Help bar
	helpBarStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Padding(0, 1)

	// Badge base style
	badgeStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)

	// Event type badges
	badgeSystem = badgeStyle.
			Background(colorMagenta).
			Foreground(lipgloss.Color("#000000"))

	badgeUser = badgeStyle.
			Background(colorYellow).
			Foreground(lipgloss.Color("#000000"))

	badgeText = badgeStyle.
			Background(colorWhite).
			Foreground(lipgloss.Color("#000000"))

	badgeTool = badgeStyle.
			Background(colorBlue).
			Foreground(lipgloss.Color("#000000"))

	badgeThinking = badgeStyle.
			Background(colorDim).
			Foreground(colorWhite)

	badgeResult = badgeStyle.
			Background(colorOrange).
			Foreground(lipgloss.Color("#000000"))

	badgeSuccess = badgeStyle.
			Background(colorGreen).
			Foreground(lipgloss.Color("#000000"))

	badgeError = badgeStyle.
			Background(colorRed).
			Foreground(colorWhite)

	badgeDefault = badgeStyle.
			Background(lipgloss.Color("#444444")).
			Foreground(colorWhite)

	// Content styles
	textStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	toolNameStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	toolInputStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	thinkingStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Italic(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(colorOrange)

	successStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	usageStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	// Event container
	eventStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			MarginBottom(1)

	// Follow mode indicator
	followOnStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	followOffStyle = lipgloss.NewStyle().
			Foreground(colorDim)
)
