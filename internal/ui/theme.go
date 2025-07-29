package ui

// Theme defines the color scheme for the UI
type Theme struct {
	Background      string
	Muted           string
	Foreground      string
	Highlight       string
	Accent          string
	Secondary       string
	Tertiary        string
	AccentBright    string
	SecondaryBright string
	TertiaryBright  string
}

// DefaultTheme provides the default color scheme
var DefaultTheme = Theme{
	Background:      "0",
	Muted:           "8",
	Foreground:      "7",
	Highlight:       "15",
	Accent:          "3",
	Secondary:       "4",
	Tertiary:        "6",
	AccentBright:    "11",
	SecondaryBright: "12",
	TertiaryBright:  "14",
}

// SetTheme updates the current color variables to match the provided theme
func SetTheme(theme Theme) {
	colorBackground = theme.Background
	colorMuted = theme.Muted
	colorForeground = theme.Foreground
	colorHighlight = theme.Highlight
	colorAccent = theme.Accent
	colorSecondary = theme.Secondary
	colorTertiary = theme.Tertiary
	colorAccentBright = theme.AccentBright
	colorSecondaryBright = theme.SecondaryBright
	colorTertiaryBright = theme.TertiaryBright
}

// Initialize with default theme
func init() {
	SetTheme(DefaultTheme)
}
