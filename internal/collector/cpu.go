package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/steam-os-monitor/monitor/pkg/metrics"
)

// CPUCollector collects CPU metrics
type CPUCollector struct{}

// NewCPUCollector creates a new CPU collector
func NewCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

// Collect gathers CPU statistics
func (c *CPUCollector) Collect() (*metrics.CPUStats, error) {
	// Get overall CPU percentage
	overallPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// Get per-core CPU percentage
	perCorePercent, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, err
	}

	// Get load average
	loadAvg, err := load.Avg()
	if err != nil {
		// Load average might not be available on all systems
		loadAvg = &load.AvgStat{}
	}

	stats := &metrics.CPUStats{
		OverallPercent: overallPercent[0],
		PerCorePercent: perCorePercent,
		LoadAvg1:       loadAvg.Load1,
		LoadAvg5:       loadAvg.Load5,
		LoadAvg15:      loadAvg.Load15,
		Timestamp:      time.Now(),
	}

	return stats, nil
}

