package ui

import "github.com/charmbracelet/lipgloss"

// Coherent color palette
const (
	// Base colors
	colorBackground = "235" // Dark gray
	colorForeground = "252" // Light gray
	colorAccent     = "108" // Green
	colorWarning    = "214" // Orange
	colorError      = "196" // Red
	colorMuted      = "240" // Dim gray
	colorHighlight  = "229" // Bright yellow
	colorInfo       = "117" // Light blue

	// Message states
	colorUnread     = "231" // Bright white
	colorRead       = "244" // Dimmed
	colorMarked     = "25"  // Blue background
	colorSelected   = "57"  // Purple background
	colorThreadLine = "238" // Thread indicator color
)

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
			BorderForeground(lipgloss.Color(colorMuted)).
			Padding(0, 1)

	styleTabActive = styleTabNormal.Border(borderTabActive, true).
			Foreground(lipgloss.Color(colorHighlight)).
			BorderForeground(lipgloss.Color(colorAccent))

	styleTabGap = styleTabNormal.BorderBottom(false).BorderLeft(false).BorderRight(false)

	styleSpinner = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccent))

	styleListTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorUnread)).
			Background(lipgloss.Color(colorMarked)).
			Padding(0, 1)

	styleListNoItems = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorMuted)).
				Align(lipgloss.Center)

	styleMessageNormal = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorForeground))

	styleMessageSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorHighlight)).
				Background(lipgloss.Color(colorSelected)).
				Bold(true)

	styleMessageMarked = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorUnread)).
				Background(lipgloss.Color(colorMarked))

	styleMessageDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted))

	styleStatusLine = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorUnread)).
			Background(lipgloss.Color(colorMarked)).
			Padding(0, 1).
			Bold(true)

	styleStatusBorder = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorSelected))

	// Message component styles - unread (bright)
	styleMsgDateUnread = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorInfo))

	styleMsgSenderUnread = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorAccent))

	styleMsgArrowUnread = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorMuted))

	styleMsgRecipientUnread = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorAccent))

	styleMsgSubjectUnread = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorUnread)).
				Bold(true)

	styleMsgTagsUnread = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorWarning))

	// Message component styles - read (dimmed)
	styleMsgDateRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorMuted))

	styleMsgSenderRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))

	styleMsgArrowRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorMuted))

	styleMsgRecipientRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))

	styleMsgSubjectRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))

	styleMsgTagsRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorMuted))

	// Thread and attachment indicators
	styleMsgThread = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorThreadLine))

	styleMsgAttachment = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorWarning))
)
