package ui

import "github.com/charmbracelet/lipgloss"

// Coherent color palette
const (
	// Base colors
	colorBackground = "0"
	colorForeground = "7"
	colorAccent     = "2"
	colorSecondary  = "12"
	colorWarning    = "11"
	colorError      = "9"
	colorMuted      = "8"
	colorHighlight  = "11"
	colorInfo       = "4"

	// Message states
	colorUnread             = "15"
	colorRead               = "8"
	colorMarkedForeground   = "0"
	colorMarkedBackground   = colorSecondary
	colorSelectedForeground = "0"
	colorSelectedBackground = colorAccent
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
			BorderForeground(lipgloss.Color(colorAccent)).
			Padding(0, 1)

	styleTabActive = styleTabNormal.Border(borderTabActive, true).
			Foreground(lipgloss.Color(colorHighlight)).
			BorderForeground(lipgloss.Color(colorAccent))

	styleTabGap = styleTabNormal.BorderBottom(false).BorderLeft(false).BorderRight(false)

	styleSpinner = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccent))

	styleListNoItems = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorMuted)).
				Align(lipgloss.Center)

	styleMessageNormal = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorForeground))

	styleStatusLine = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorBackground)).
			Background(lipgloss.Color(colorAccent)).
			Padding(0, 1).
			Bold(true)

	stylePaginationActivePage = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorSecondary))

	stylePaginationInactivePage = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorAccent))

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
				Foreground(lipgloss.Color(colorRead))

	styleMsgSenderRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))

	styleMsgArrowRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))

	styleMsgRecipientRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))

	styleMsgSubjectRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))

	styleMsgTagsRead = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorRead))
)
