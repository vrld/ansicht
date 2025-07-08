package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/service"
	"github.com/vrld/ansicht/internal/ui"
)

func main() {
	messages := service.NewMessages()
	queries := service.NewQueries()

	// Load configuration
	runtime, err := runtime.LoadRuntime()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	runtime.Messages = messages

	// Initialize UI model
	model := ui.NewModel(messages, queries)
	model.KeyReceiver = runtime

	// Start the application
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}
