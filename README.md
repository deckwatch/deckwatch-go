# SteamOS System Monitoring & Diagnostic Tool

A comprehensive system monitoring tool for SteamOS/SteamDeck built with Go, featuring a desktop GUI with btop-inspired visualizations.

## Features

- **Real-time System Monitoring**
  - CPU usage (per-core and overall)
  - Memory usage (RAM and Swap)
  - Disk I/O and usage statistics
  - Network bandwidth and packet statistics
  - Game performance metrics (FPS, frame times)
  - Steam-specific metrics (download speeds, library status)

- **Beautiful GUI**
  - btop-inspired visualizations
  - Color-coded metrics (green/yellow/red based on usage)
  - Real-time updates
  - Scrollable interface

- **Comprehensive Logging**
  - Separate log files for each metric type
  - JSON or CSV format support
  - Configurable log directory

## Requirements

- Go 1.21 or later
- SteamOS/SteamDeck or Linux system
- Fyne GUI framework dependencies

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd steam-os-monitor
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o steam-os-monitor ./cmd/monitor
```

## Usage

Run the application:
```bash
./steam-os-monitor
```

Or with a custom config file:
```bash
./steam-os-monitor -config /path/to/config.yaml
```

## Configuration

The application creates a default configuration file at `~/.steam-os-monitor/config.yaml` on first run. You can customize:

- Refresh rate (milliseconds)
- Log directory and format
- Widget visibility
- Color themes
- Steam API credentials

Example configuration:
```yaml
refresh_rate: 1000
log_dir: ~/.steam-os-monitor/logs
log_format: json
widgets:
  show_cpu: true
  show_memory: true
  show_disk: true
  show_network: true
  show_game: true
  show_steam: true
theme:
  background_color: "#1e1e2e"
  text_color: "#cdd6f4"
  bar_color: "#89b4fa"
  bar_color_high: "#f38ba8"
  bar_color_medium: "#fab387"
  bar_color_low: "#a6e3a1"
```

## Log Files

Logs are stored in separate files:
- `cpu.log` - CPU metrics
- `memory.log` - Memory metrics
- `disk.log` - Disk metrics
- `network.log` - Network metrics
- `game_performance.log` - Game performance metrics
- `steam.log` - Steam metrics

## Project Structure

```
steam-os-monitor/
├── cmd/monitor/          # Application entry point
├── internal/
│   ├── collector/        # System metrics collection
│   ├── logger/           # Logging functionality
│   ├── ui/               # GUI components
│   └── config/           # Configuration management
└── pkg/metrics/          # Metric data structures
```

## Development

To contribute or modify:

1. Make your changes
2. Run tests (if available):
```bash
go test ./...
```

3. Build and test:
```bash
go build -o steam-os-monitor ./cmd/monitor
```

## License

[Add your license here]

## Notes

- Game performance metrics require gamescope integration (may need additional setup)
- Steam metrics require Steam API access (may need API key configuration)
- Some metrics may require elevated permissions on certain systems

