package main

import (
	"fmt"
	"os"

	"github.com/vibolsovichea/scripture/internal/config"
	"github.com/vibolsovichea/scripture/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cfg == nil {
		// First run — prompt for vault location
		p := tea.NewProgram(ui.NewSetupModel())
		m, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		setup := m.(ui.SetupModel)
		if setup.Cancelled {
			fmt.Println("Setup cancelled.")
			os.Exit(0)
		}

		cfg, err = config.Init(setup.VaultPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing vault: %v\n", err)
			os.Exit(1)
		}
	}

	p := tea.NewProgram(ui.NewAppModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
