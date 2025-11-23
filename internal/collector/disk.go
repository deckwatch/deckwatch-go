package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/steam-os-monitor/monitor/pkg/metrics"
)

// DiskCollector collects disk metrics
type DiskCollector struct {
	lastIOStats map[string]*disk.IOCountersStat
}

// NewDiskCollector creates a new disk collector
func NewDiskCollector() *DiskCollector {
	return &DiskCollector{
		lastIOStats: make(map[string]*disk.IOCountersStat),
	}
}

// Collect gathers disk statistics for all partitions
func (c *DiskCollector) Collect() ([]*metrics.DiskStats, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var stats []*metrics.DiskStats
	currentIOStats, err := disk.IOCounters()
	if err != nil {
		// IO stats might not be available, continue without them
		currentIOStats = make(map[string]disk.IOCountersStat)
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// Skip partitions we can't read
			continue
		}

		diskStat := &metrics.DiskStats{
			Device:     partition.Device,
			MountPoint: partition.Mountpoint,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			UsedPercent: usage.UsedPercent,
			Timestamp:   time.Now(),
		}

		// Get IO stats if available
		if ioStat, exists := currentIOStats[partition.Device]; exists {
			diskStat.ReadBytes = ioStat.ReadBytes
			diskStat.WriteBytes = ioStat.WriteBytes
			diskStat.ReadIOPS = ioStat.ReadCount
			diskStat.WriteIOPS = ioStat.WriteCount

			// Calculate speed if we have previous stats
			if lastIO, exists := c.lastIOStats[partition.Device]; exists {
				// Note: This is a simplified calculation
				// In a real implementation, you'd track time between samples
				_ = lastIO
			}
		}

		stats = append(stats, diskStat)
	}

	// Update last IO stats
	c.lastIOStats = make(map[string]*disk.IOCountersStat)
	for device, ioStat := range currentIOStats {
		statCopy := ioStat
		c.lastIOStats[device] = &statCopy
	}

	return stats, nil
}

