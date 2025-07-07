package ui

import "github.com/charmbracelet/lipgloss"

// styleTabNormal defines the styling for tabs
var (
	borderTabActive = lipgloss.Border{
		Top:         " ",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopRight:    "┌",
		TopLeft:     "┐",
		BottomRight: "╯",
		BottomLeft:  "╰",
	}

	borderTab = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopRight:    "┬",
		TopLeft:     "┬",
		BottomRight: "╯",
		BottomLeft:  "╰",
	}

	styleTabNormal = lipgloss.NewStyle().
			Border(borderTab, true).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	styleTabActive = styleTabNormal.Border(borderTabActive, true)

	styleTabGap = styleTabNormal.BorderBottom(false).BorderLeft(false).BorderRight(false)

	styleSpinner = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62"))

	styleListTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("25")).
			Padding(0, 1)

	styleListNoItems = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Align(lipgloss.Center)

	styleMessageNormal = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	styleMessageSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(true)

	styleMessageMarked = lipgloss.NewStyle().
				Foreground(lipgloss.Color("231")).
				Background(lipgloss.Color("25"))

	styleMessageDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
