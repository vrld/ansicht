package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/ui"
)

func main() {
	// Load configuration
	runtime, err := runtime.LoadRuntime()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize UI model
	model := ui.NewModel(runtime)

	// Start the application
	p := tea.NewProgram(model, tea.WithAltScreen())
	runtime.Controller = &ui.RuntimeAdapter{Program: p}
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}
