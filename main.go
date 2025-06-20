package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/ui"
)

func main() {
	model := ui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}
