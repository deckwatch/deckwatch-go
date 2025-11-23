package collector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/steam-os-monitor/monitor/pkg/metrics"
)

// SteamCollector collects Steam-specific metrics
type SteamCollector struct {
	steamDir string
}

// NewSteamCollector creates a new Steam collector
func NewSteamCollector() *SteamCollector {
	steamDir := os.Getenv("HOME")
	if steamDir == "" {
		steamDir = "/home/deck"
	}
	steamDir = filepath.Join(steamDir, ".steam", "steam")

	return &SteamCollector{
		steamDir: steamDir,
	}
}

// Collect gathers Steam statistics
func (c *SteamCollector) Collect() (*metrics.SteamStats, error) {
	stats := &metrics.SteamStats{
		DownloadProgress: make(map[string]float64),
		Timestamp:        time.Now(),
	}

	// Get download speed from Steam
	downloadSpeed, err := c.getDownloadSpeed()
	if err == nil {
		stats.DownloadSpeed = downloadSpeed
	}

	// Get upload speed
	uploadSpeed, err := c.getUploadSpeed()
	if err == nil {
		stats.UploadSpeed = uploadSpeed
	}

	// Get active downloads count
	activeDownloads, err := c.getActiveDownloads()
	if err == nil {
		stats.ActiveDownloads = activeDownloads
	}

	// Get library size
	librarySize, installedGames, err := c.getLibraryInfo()
	if err == nil {
		stats.LibrarySize = librarySize
		stats.InstalledGames = installedGames
	}

	// Get download progress for active downloads
	progress, err := c.getDownloadProgress()
	if err == nil {
		stats.DownloadProgress = progress
	}

	return stats, nil
}

// getDownloadSpeed attempts to get Steam download speed
func (c *SteamCollector) getDownloadSpeed() (float64, error) {
	// Try to read from Steam's download stats
	// This is a placeholder - actual implementation would query Steam API or read Steam logs

	// Check for Steam download stats file
	downloadStatsFile := filepath.Join(c.steamDir, "logs", "download_stats.txt")
	if data, err := os.ReadFile(downloadStatsFile); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.Contains(line, "speed") || strings.Contains(line, "bytes/sec") {
				parts := strings.Fields(line)
				for i, part := range parts {
					if strings.Contains(part, "bytes/sec") || strings.Contains(part, "B/s") {
						if i > 0 {
							if speed, err := strconv.ParseFloat(parts[i-1], 64); err == nil {
								return speed, nil
							}
						}
					}
				}
			}
		}
	}

	// Try to get from network interface if Steam is downloading
	// This is a simplified approach
	return 0, fmt.Errorf("could not get download speed")
}

// getUploadSpeed attempts to get Steam upload speed
func (c *SteamCollector) getUploadSpeed() (float64, error) {
	// Similar to download speed
	return 0, fmt.Errorf("could not get upload speed")
}

// getActiveDownloads gets the count of active downloads
func (c *SteamCollector) getActiveDownloads() (int, error) {
	// Try to query Steam process or API
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count := 0
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "steam") && strings.Contains(line, "download") {
			count++
		}
	}

	return count, nil
}

// getLibraryInfo gets library size and installed games count
func (c *SteamCollector) getLibraryInfo() (uint64, int, error) {
	libraryPath := filepath.Join(c.steamDir, "steamapps", "common")
	
	var totalSize uint64
	var gameCount int

	err := filepath.Walk(libraryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() {
			// Count game directories
			parent := filepath.Dir(path)
			if parent == libraryPath {
				gameCount++
			}
		} else {
			totalSize += uint64(info.Size())
		}
		return nil
	})

	if err != nil {
		return 0, 0, err
	}

	return totalSize, gameCount, nil
}

// getDownloadProgress gets download progress for active downloads
func (c *SteamCollector) getDownloadProgress() (map[string]float64, error) {
	progress := make(map[string]float64)
	
	// This would require Steam API integration
	// For now, return empty map
	return progress, nil
}

