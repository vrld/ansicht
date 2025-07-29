package runtime

import (
	lua "github.com/Shopify/go-lua"
)

// ThemeData represents theme color values to pass to UI without import cycle
type ThemeData struct {
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

// luaSetTheme implements ansicht.set_theme() function
// Expects a Lua table with theme color fields
func (r *Runtime) luaSetTheme(L *lua.State) int {
	if !L.IsTable(1) {
		lua.Errorf(L, "set_theme expects a table")
		panic("unreachable")
	}

	// Parse theme fields from Lua table with defaults
	theme := ThemeData{
		Background:      lFieldStringOrDefault(L, 1, "background", "0"),
		Muted:           lFieldStringOrDefault(L, 1, "muted", "8"),
		Foreground:      lFieldStringOrDefault(L, 1, "foreground", "7"),
		Highlight:       lFieldStringOrDefault(L, 1, "highlight", "15"),
		Accent:          lFieldStringOrDefault(L, 1, "accent", "3"),
		Secondary:       lFieldStringOrDefault(L, 1, "secondary", "4"),
		Tertiary:        lFieldStringOrDefault(L, 1, "tertiary", "6"),
		AccentBright:    lFieldStringOrDefault(L, 1, "accent_bright", "11"),
		SecondaryBright: lFieldStringOrDefault(L, 1, "secondary_bright", "12"),
		TertiaryBright:  lFieldStringOrDefault(L, 1, "tertiary_bright", "14"),
	}

	r.Controller.SetTheme(theme)
	return 0
}
