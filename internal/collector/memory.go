package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/steam-os-monitor/monitor/pkg/metrics"
)

// MemoryCollector collects memory metrics
type MemoryCollector struct{}

// NewMemoryCollector creates a new memory collector
func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

// Collect gathers memory statistics
func (c *MemoryCollector) Collect() (*metrics.MemoryStats, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	stats := &metrics.MemoryStats{
		Total:      vmStat.Total,
		Used:       vmStat.Used,
		Available:  vmStat.Available,
		UsedPercent: vmStat.UsedPercent,
		SwapTotal:   swapStat.Total,
		SwapUsed:    swapStat.Used,
		SwapPercent: swapStat.UsedPercent,
		Timestamp:   time.Now(),
	}

	return stats, nil
}

