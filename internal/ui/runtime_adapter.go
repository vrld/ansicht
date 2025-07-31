package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/service"
)

type RuntimeAdapter struct {
	Program *tea.Program
}

func (a *RuntimeAdapter) Quit() {
	go a.Program.Send(tea.QuitMsg{})
}

func (a *RuntimeAdapter) Refresh() {
	go a.Program.Send(Refresh{})
}

func (a *RuntimeAdapter) Status(message string) {
	service.Status().Set(message)
	go a.Program.Send(message) // unhandled, but causes a redraw
}

func (a *RuntimeAdapter) Notify(message string, level string, timeout int) {
	go a.Program.Send(NotifyMsg{
		Message: message,
		Level:   level,
		Timeout: timeout,
	})
}

func (a *RuntimeAdapter) Input(prompt, placeholder string) {
	go a.Program.Send(OpenInputEvent{
		Placeholder: placeholder,
		Prompt:      prompt,
	})
}

func (a *RuntimeAdapter) SpawnResult(result runtime.SpawnResult) {
	go a.Program.Send(result)
}

func (a *RuntimeAdapter) QueryNew(query string) {
	go a.Program.Send(QueryNewMsg{query})
}

func (a *RuntimeAdapter) QuerySelectNext() {
	go a.Program.Send(QueryNextMsg{})
}

func (a *RuntimeAdapter) QuerySelectPrev() {
	go a.Program.Send(QueryPrevMsg{})
}

func (a *RuntimeAdapter) MarksToggle() {
	go a.Program.Send(MarksToggleMsg{})
}

func (a *RuntimeAdapter) MarksInvert() {
	go a.Program.Send(MarksInvertMsg{})
}

func (a *RuntimeAdapter) MarksClear() {
	go a.Program.Send(MarksClearMsg{})
}

func (a *RuntimeAdapter) SetTheme(theme any) {
	if theme, ok := theme.(runtime.ThemeData); ok {
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
		go a.Program.Send("theme_updated") // trigger redraw
	}
}
