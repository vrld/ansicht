package runtime

import lua "github.com/Shopify/go-lua"

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
	Warning         string
	Error           string
}

// ansicht.theme.set{ ... }
func (r *Runtime) luaThemeSet(L *lua.State) int {
	if !L.IsTable(1) {
		lua.Errorf(L, "set_theme expects a table")
		panic("unreachable")
	}

	r.Controller.SetTheme(ThemeData{
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
		Warning:         lFieldStringOrDefault(L, 1, "warning", "13"),
		Error:           lFieldStringOrDefault(L, 1, "error", "9"),
	})
	return 0
}
