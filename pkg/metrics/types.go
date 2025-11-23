package metrics

import "time"

// CPUStats represents CPU usage metrics
type CPUStats struct {
	OverallPercent float64   `json:"overall_percent"`
	PerCorePercent []float64 `json:"per_core_percent"`
	LoadAvg1       float64   `json:"load_avg_1"`
	LoadAvg5       float64   `json:"load_avg_5"`
	LoadAvg15      float64   `json:"load_avg_15"`
	Timestamp      time.Time `json:"timestamp"`
}

// MemoryStats represents memory usage metrics
type MemoryStats struct {
	Total       uint64    `json:"total"`
	Used        uint64    `json:"used"`
	Available   uint64    `json:"available"`
	UsedPercent  float64   `json:"used_percent"`
	SwapTotal    uint64    `json:"swap_total"`
	SwapUsed     uint64    `json:"swap_used"`
	SwapPercent  float64   `json:"swap_percent"`
	Timestamp    time.Time `json:"timestamp"`
}

// DiskStats represents disk I/O metrics
type DiskStats struct {
	Device      string    `json:"device"`
	MountPoint  string    `json:"mount_point"`
	Total       uint64    `json:"total"`
	Used        uint64    `json:"used"`
	Free        uint64    `json:"free"`
	UsedPercent float64   `json:"used_percent"`
	ReadBytes   uint64    `json:"read_bytes"`
	WriteBytes  uint64    `json:"write_bytes"`
	ReadIOPS    uint64    `json:"read_iops"`
	WriteIOPS   uint64    `json:"write_iops"`
	Timestamp   time.Time `json:"timestamp"`
}

// NetworkStats represents network statistics
type NetworkStats struct {
	Interface    string    `json:"interface"`
	BytesSent    uint64    `json:"bytes_sent"`
	BytesRecv    uint64    `json:"bytes_recv"`
	PacketsSent  uint64    `json:"packets_sent"`
	PacketsRecv  uint64    `json:"packets_recv"`
	SpeedSent    float64   `json:"speed_sent"`    // bytes per second
	SpeedRecv    float64   `json:"speed_recv"`    // bytes per second
	Timestamp    time.Time `json:"timestamp"`
}

// GamePerformanceStats represents game performance metrics
type GamePerformanceStats struct {
	FPS         float64   `json:"fps"`
	FrameTime   float64   `json:"frame_time_ms"`
	FrameTimeMin float64  `json:"frame_time_min_ms"`
	FrameTimeMax float64  `json:"frame_time_max_ms"`
	GameName    string    `json:"game_name"`
	Timestamp   time.Time `json:"timestamp"`
}

// SteamStats represents Steam-specific metrics
type SteamStats struct {
	DownloadSpeed    float64   `json:"download_speed_bytes_per_sec"`
	UploadSpeed      float64   `json:"upload_speed_bytes_per_sec"`
	ActiveDownloads  int       `json:"active_downloads"`
	LibrarySize      uint64    `json:"library_size_bytes"`
	InstalledGames   int       `json:"installed_games"`
	DownloadProgress map[string]float64 `json:"download_progress"` // game_id -> progress percentage
	Timestamp        time.Time `json:"timestamp"`
}

