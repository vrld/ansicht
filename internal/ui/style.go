package ui

import "github.com/charmbracelet/lipgloss"

// Coherent color palette
const (
	// Base colors
	colorBackground = "0"
	colorMuted      = "8"
	colorForeground = "7"
	colorHighlight  = "15"

	colorAccent    = "3"
	colorSecondary = "4"
	colorTertiary  = "6"

	colorAccentBright    = "11"
	colorSecondaryBright = "12"
	colorTertiaryBright  = "14"
)

// styleTabNormal defines the styling for tabs
var (
	styleSpinner = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccent))

	styleListNoItems = lipgloss.NewStyle().
				Bold(true).
				Align(lipgloss.Center, lipgloss.Center)

	styleStatusLine = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorBackground)).
			Background(lipgloss.Color(colorSecondary)).
			Padding(0, 1).
			Bold(true)
)
