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
	inputHistory := service.NewInputHistory()
	status := service.NewStatus()

	// Load configuration
	runtime, err := runtime.LoadRuntime()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	runtime.Messages = messages
	runtime.Status = status

	// Initialize UI model
	model := ui.NewModel()
	model.Messages = messages
	model.Queries = queries
	model.InputHistory = inputHistory
	model.Status = status
	model.Runtime = runtime

	// Start the application
	p := tea.NewProgram(model, tea.WithAltScreen())
	runtime.SendMessage = p.Send
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}
