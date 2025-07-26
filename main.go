package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/ansicht/internal/runtime"
	"github.com/vrld/ansicht/internal/service"
	"github.com/vrld/ansicht/internal/ui"
)

func main() {
	handleCommandline()
	defer service.Logger().Close()

	runtime, err := runtime.LoadRuntime()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	model := ui.NewModel(runtime)

	p := tea.NewProgram(model, tea.WithAltScreen())
	runtime.Controller = &ui.RuntimeAdapter{Program: p}
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}

func handleCommandline() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "ansicht - email at a glance\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", filepath.Base(os.Args[0]))

		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()

		fmt.Fprintf(flag.CommandLine.Output(), "\nConfiguration:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  ansicht looks for configuration in:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    $XDG_CONFIG_HOME/ansicht/init.lua\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    ~/.config/ansicht/init.lua\n")
	}

	logFile := flag.String("log-file", "", "Write logs to this file if given")
	help := flag.Bool("h", false, "Show help message")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *logFile != "" {
		if err := service.Logger().Initialize(*logFile); err != nil {
			log.Fatalf("Error initializing logging: %v", err)
		}
	}
}
