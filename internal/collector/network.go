package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/steam-os-monitor/monitor/pkg/metrics"
)

// NetworkCollector collects network metrics
type NetworkCollector struct {
	lastStats map[string]*net.IOCountersStat
	lastTime  time.Time
}

// NewNetworkCollector creates a new network collector
func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{
		lastStats: make(map[string]*net.IOCountersStat),
		lastTime:  time.Now(),
	}
}

// Collect gathers network statistics
func (c *NetworkCollector) Collect() ([]*metrics.NetworkStats, error) {
	currentStats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	elapsed := now.Sub(c.lastTime).Seconds()
	if elapsed == 0 {
		elapsed = 1.0 // Avoid division by zero
	}

	var stats []*metrics.NetworkStats

	for _, stat := range currentStats {
		// Skip loopback interfaces
		if stat.Name == "lo" || stat.Name == "lo0" {
			continue
		}

		networkStat := &metrics.NetworkStats{
			Interface:   stat.Name,
			BytesSent:   stat.BytesSent,
			BytesRecv:   stat.BytesRecv,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
			Timestamp:   now,
		}

		// Calculate speed if we have previous stats
		if lastStat, exists := c.lastStats[stat.Name]; exists {
			bytesSentDiff := float64(stat.BytesSent - lastStat.BytesSent)
			bytesRecvDiff := float64(stat.BytesRecv - lastStat.BytesRecv)

			networkStat.SpeedSent = bytesSentDiff / elapsed
			networkStat.SpeedRecv = bytesRecvDiff / elapsed
		}

		stats = append(stats, networkStat)
	}

	// Update last stats
	c.lastStats = make(map[string]*net.IOCountersStat)
	for _, stat := range currentStats {
		statCopy := stat
		c.lastStats[stat.Name] = &statCopy
	}
	c.lastTime = now

	return stats, nil
}

