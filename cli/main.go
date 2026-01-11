package main

import (
	// "bufio"
	"context"
	"fmt"
	"flag"
	// "log"
	"os"
	"os/signal"
	// "strings"
	"syscall"
	// "time"

	"carmaintenance/internal/core"
	"carmaintenance/cli/cliapp"
)

func main() {
	dbPath := flag.String("dbpath", "$HOME/carmaintenance/database.db", "the path of the database")
	specsDir := flag.String("specs", "$HOME/carmaintenance/specs/", "the directory with the specifications YAML files")
	flag.Parse()

	// Create context that cancels on interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// backend, err := core.NewBackend(dbPath, specsDir)
	backend, err := core.NewBackend(dbPath, specsDir)
	if err != nil {
		fmt.Printf("Error: %w\n", err)
		os.Exit(1)
	}
	defer backend.Cleanup()

	errCh := make(chan error, 1)

	// Run the main CLI app logic
	go func() {
		// TODO: Initialize Bubbletea

		if err := cliapp.Run(ctx, backend); err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Shutting down the app")
	case err := <-errCh:
		if err != nil {
			fmt.Printf("CLI app error: %w\n", err)
		} else {
			// TODO: Remove this case
			fmt.Println("CLI app exited normally")
		}
	}
}
