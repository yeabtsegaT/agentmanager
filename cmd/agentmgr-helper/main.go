// AgentManager Helper - System tray helper for managing AI development agents
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kevinelliott/agentmgr/internal/systray"
	"github.com/kevinelliott/agentmgr/pkg/catalog"
	"github.com/kevinelliott/agentmgr/pkg/config"
	"github.com/kevinelliott/agentmgr/pkg/detector"
	"github.com/kevinelliott/agentmgr/pkg/installer"
	"github.com/kevinelliott/agentmgr/pkg/platform"
	"github.com/kevinelliott/agentmgr/pkg/storage"
)

// Version information (set by build flags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Get current platform
	plat := platform.Current()

	// Load configuration
	loader := config.NewLoader()
	cfg, err := loader.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize storage
	dataDir := plat.GetDataDir()
	store, err := storage.NewSQLiteStore(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Initialize the database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := store.Initialize(ctx); err != nil {
		cancel()
		fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
	cancel()

	// Initialize detector with strategies
	det := detector.New(plat)

	// Initialize catalog manager
	cat := catalog.NewManager(cfg, store)

	// Initialize installer manager
	inst := installer.NewManager(plat)

	// Create systray app
	app := systray.New(cfg, loader, plat, store, det, cat, inst, version)

	// Handle shutdown signals in a goroutine
	// (systray.Run must be on main thread for macOS)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
		app.Quit()
	}()

	// Run systray on main thread (required for macOS)
	// This blocks until systray.Quit() is called
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
