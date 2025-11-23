package collector

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/steam-os-monitor/monitor/pkg/metrics"
)

// GameCollector collects game performance metrics
// This integrates with gamescope and Steam APIs to get FPS and frame time data
type GameCollector struct {
	lastFrameTime time.Time
}

// NewGameCollector creates a new game collector
func NewGameCollector() *GameCollector {
	return &GameCollector{}
}

// Collect gathers game performance statistics
func (c *GameCollector) Collect() (*metrics.GamePerformanceStats, error) {
	stats := &metrics.GamePerformanceStats{
		Timestamp: time.Now(),
	}

	// Try to get FPS from gamescope
	fps, err := c.getFPSFromGamescope()
	if err == nil {
		stats.FPS = fps
		// Calculate frame time from FPS
		if fps > 0 {
			stats.FrameTime = 1000.0 / fps // milliseconds
		}
	}

	// Try to get game name from Steam
	gameName, err := c.getCurrentGameName()
	if err == nil {
		stats.GameName = gameName
	}

	// Frame time min/max would require tracking over time
	// For now, set them to the current frame time
	if stats.FrameTime > 0 {
		stats.FrameTimeMin = stats.FrameTime
		stats.FrameTimeMax = stats.FrameTime
	}

	return stats, nil
}

// getFPSFromGamescope attempts to get FPS from gamescope
// This is a placeholder - actual implementation would need gamescope integration
func (c *GameCollector) getFPSFromGamescope() (float64, error) {
	// Try to read from gamescope's FPS counter
	// On SteamOS, gamescope might expose FPS via environment variables or files
	// This is a simplified implementation

	// Check for gamescope FPS file (common location)
	fpsFile := "/tmp/gamescope-fps"
	if data, err := os.ReadFile(fpsFile); err == nil {
		if fps, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
			return fps, nil
		}
	}

	// Try to get FPS from gamescope command if available
	cmd := exec.Command("gamescope", "--fps")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if strings.Contains(line, "fps") {
				parts := strings.Fields(line)
				for i, part := range parts {
					if strings.Contains(part, "fps") && i > 0 {
						if fps, err := strconv.ParseFloat(parts[i-1], 64); err == nil {
							return fps, nil
						}
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("could not get FPS from gamescope")
}

// getCurrentGameName attempts to get the current game name
func (c *GameCollector) getCurrentGameName() (string, error) {
	// Try to get from Steam process or environment
	// This is a placeholder - actual implementation would query Steam API

	// Check for Steam game name in environment
	if gameName := os.Getenv("STEAM_GAME_NAME"); gameName != "" {
		return gameName, nil
	}

	// Try to get from Steam process
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "steam") && strings.Contains(line, "game") {
				// Extract game name from process line
				parts := strings.Fields(line)
				for _, part := range parts {
					if strings.Contains(part, ".exe") || strings.Contains(part, "game") {
						return part, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("could not determine game name")
}

