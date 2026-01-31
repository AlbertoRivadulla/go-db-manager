package main

import (
	// "bufio"
	"context"
	"fmt"
	"flag"
	// "log"
	"os"
	"os/signal"
	"syscall"

	"carmaintenance/internal/core"
	"carmaintenance/cli/cliapp"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dbPath := flag.String("dbpath", "$HOME/carmaintenance/database.db", "the path of the database")
	specsDir := flag.String("specs", "$HOME/carmaintenance/specs/", "the directory with the specifications YAML files")
	flag.Parse()

	// Create context that cancels on interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	backend, err := core.NewBackend(dbPath, specsDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer backend.Cleanup()

	errCh := make(chan error, 1)

	// Run the main CLI app logic
	go func() {
		// TODO: Initialize Bubbletea
		model := cliapp.NewModel(ctx, backend)
		p := tea.NewProgram(model)
		// p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return
	case err := <-errCh:
		if err != nil {
			fmt.Printf("CLI app error: %w\n", err)
		}
	}
}
