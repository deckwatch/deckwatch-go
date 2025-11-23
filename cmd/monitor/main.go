package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/steam-os-monitor/monitor/internal/config"
	"github.com/steam-os-monitor/monitor/internal/logger"
	"github.com/steam-os-monitor/monitor/internal/ui"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", getDefaultConfigPath(), "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.LogDir, cfg.LogFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Create and show window
	window, err := ui.NewWindow(cfg, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating window: %v\n", err)
		os.Exit(1)
	}

	// Run application
	window.ShowAndRun()
}

func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./config.yaml"
	}
	return filepath.Join(homeDir, ".steam-os-monitor", "config.yaml")
}

