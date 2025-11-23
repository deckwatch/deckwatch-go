package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	RefreshRate int      `yaml:"refresh_rate"` // milliseconds
	LogDir      string   `yaml:"log_dir"`
	LogFormat   string   `yaml:"log_format"` // "json" or "csv"
	Widgets     Widgets `yaml:"widgets"`
	Theme       Theme    `yaml:"theme"`
	Steam       Steam    `yaml:"steam"`
}

// Widgets configuration
type Widgets struct {
	ShowCPU     bool `yaml:"show_cpu"`
	ShowMemory  bool `yaml:"show_memory"`
	ShowDisk    bool `yaml:"show_disk"`
	ShowNetwork bool `yaml:"show_network"`
	ShowGame    bool `yaml:"show_game"`
	ShowSteam   bool `yaml:"show_steam"`
}

// Theme configuration
type Theme struct {
	BackgroundColor string `yaml:"background_color"`
	TextColor       string `yaml:"text_color"`
	BarColor        string `yaml:"bar_color"`
	BarColorHigh    string `yaml:"bar_color_high"`
	BarColorMedium  string `yaml:"bar_color_medium"`
	BarColorLow     string `yaml:"bar_color_low"`
}

// Steam configuration
type Steam struct {
	APIKey string `yaml:"api_key"`
}

// LoadConfig loads configuration from file or creates default
func LoadConfig(configPath string) (*Config, error) {
	// Default configuration
	defaultConfig := &Config{
		RefreshRate: 1000, // 1 second
		LogDir:      getDefaultLogDir(),
		LogFormat:   "json",
		Widgets: Widgets{
			ShowCPU:     true,
			ShowMemory:  true,
			ShowDisk:    true,
			ShowNetwork: true,
			ShowGame:    true,
			ShowSteam:   true,
		},
		Theme: Theme{
			BackgroundColor: "#1e1e2e",
			TextColor:       "#cdd6f4",
			BarColor:        "#89b4fa",
			BarColorHigh:    "#f38ba8",
			BarColorMedium:  "#fab387",
			BarColorLow:     "#a6e3a1",
		},
	}

	// Try to load from file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		if err := SaveConfig(configPath, defaultConfig); err != nil {
			return defaultConfig, fmt.Errorf("failed to create default config: %w", err)
		}
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return defaultConfig, fmt.Errorf("failed to parse config: %w", err)
	}

	// Merge with defaults for missing fields
	if config.RefreshRate == 0 {
		config.RefreshRate = defaultConfig.RefreshRate
	}
	if config.LogDir == "" {
		config.LogDir = defaultConfig.LogDir
	}
	if config.LogFormat == "" {
		config.LogFormat = defaultConfig.LogFormat
	}

	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(configPath string, config *Config) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func getDefaultLogDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./logs"
	}
	return filepath.Join(homeDir, ".steam-os-monitor", "logs")
}

