package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/steam-os-monitor/monitor/internal/collector"
	"github.com/steam-os-monitor/monitor/internal/config"
	"github.com/steam-os-monitor/monitor/internal/logger"
	"github.com/steam-os-monitor/monitor/internal/theme"
	"github.com/steam-os-monitor/monitor/internal/ui/widgets"
)

// Window represents the main application window
type Window struct {
	app        fyne.App
	window     fyne.Window
	config     *config.Config
	logger     *logger.Logger
	theme      *theme.Theme
	
	// Collectors
	cpuCollector     *collector.CPUCollector
	memoryCollector  *collector.MemoryCollector
	diskCollector    *collector.DiskCollector
	networkCollector *collector.NetworkCollector
	gameCollector    *collector.GameCollector
	steamCollector   *collector.SteamCollector
	
	// Widgets
	cpuWidget     *widgets.CPUWidget
	memoryWidget *widgets.MemoryWidget
	diskWidget   *widgets.DiskWidget
	networkWidget *widgets.NetworkWidget
	gameWidget   *widgets.GameWidget
	steamWidget  *widgets.SteamWidget
	
	// Container
	content *container.Scroll
	
	// Update ticker
	ticker *time.Ticker
}

// NewWindow creates a new application window
func NewWindow(cfg *config.Config, log *logger.Logger) (*Window, error) {
	application := app.NewWithID("steam-os-monitor")
	
	w := &Window{
		app:    application,
		config: cfg,
		logger: log,
		theme:  theme.DefaultTheme(),
	}
	
	// Initialize collectors
	w.cpuCollector = collector.NewCPUCollector()
	w.memoryCollector = collector.NewMemoryCollector()
	w.diskCollector = collector.NewDiskCollector()
	w.networkCollector = collector.NewNetworkCollector()
	w.gameCollector = collector.NewGameCollector()
	w.steamCollector = collector.NewSteamCollector()
	
	// Apply theme
	theme.ApplyTheme(application, w.theme)
	
	// Create window
	w.window = application.NewWindow("SteamOS System Monitor")
	w.window.Resize(fyne.NewSize(1200, 800))
	w.window.CenterOnScreen()
	
	// Create widgets
	w.cpuWidget = widgets.NewCPUWidget(w.theme)
	w.memoryWidget = widgets.NewMemoryWidget(w.theme)
	w.diskWidget = widgets.NewDiskWidget(w.theme)
	w.networkWidget = widgets.NewNetworkWidget(w.theme)
	w.gameWidget = widgets.NewGameWidget(w.theme)
	w.steamWidget = widgets.NewSteamWidget(w.theme)
	
	// Create layout
	w.setupLayout()
	
	// Setup update loop
	w.setupUpdateLoop()
	
	return w, nil
}

// setupLayout creates the window layout
func (w *Window) setupLayout() {
	var widgetContainers []fyne.CanvasObject
	
	if w.config.Widgets.ShowCPU {
		widgetContainers = append(widgetContainers, w.cpuWidget)
	}
	if w.config.Widgets.ShowMemory {
		widgetContainers = append(widgetContainers, w.memoryWidget)
	}
	if w.config.Widgets.ShowDisk {
		widgetContainers = append(widgetContainers, w.diskWidget)
	}
	if w.config.Widgets.ShowNetwork {
		widgetContainers = append(widgetContainers, w.networkWidget)
	}
	if w.config.Widgets.ShowGame {
		widgetContainers = append(widgetContainers, w.gameWidget)
	}
	if w.config.Widgets.ShowSteam {
		widgetContainers = append(widgetContainers, w.steamWidget)
	}
	
	// Create scrollable container with grid layout
	content := container.NewVBox(widgetContainers...)
	w.content = container.NewScroll(content)
	w.window.SetContent(w.content)
}

// setupUpdateLoop sets up the periodic update loop
func (w *Window) setupUpdateLoop() {
	refreshDuration := time.Duration(w.config.RefreshRate) * time.Millisecond
	w.ticker = time.NewTicker(refreshDuration)
	
	go func() {
		for range w.ticker.C {
			w.updateMetrics()
		}
	}()
	
	// Initial update
	w.updateMetrics()
}

// updateMetrics collects and updates all metrics
func (w *Window) updateMetrics() {
	// Collect CPU metrics
	if w.config.Widgets.ShowCPU {
		cpuStats, err := w.cpuCollector.Collect()
		if err == nil {
			w.cpuWidget.Update(cpuStats)
			w.logger.LogCPU(cpuStats)
		}
	}
	
	// Collect Memory metrics
	if w.config.Widgets.ShowMemory {
		memStats, err := w.memoryCollector.Collect()
		if err == nil {
			w.memoryWidget.Update(memStats)
			w.logger.LogMemory(memStats)
		}
	}
	
	// Collect Disk metrics
	if w.config.Widgets.ShowDisk {
		diskStats, err := w.diskCollector.Collect()
		if err == nil {
			w.diskWidget.Update(diskStats)
			for _, stat := range diskStats {
				w.logger.LogDisk(stat)
			}
		}
	}
	
	// Collect Network metrics
	if w.config.Widgets.ShowNetwork {
		netStats, err := w.networkCollector.Collect()
		if err == nil {
			w.networkWidget.Update(netStats)
			for _, stat := range netStats {
				w.logger.LogNetwork(stat)
			}
		}
	}
	
	// Collect Game metrics
	if w.config.Widgets.ShowGame {
		gameStats, err := w.gameCollector.Collect()
		if err == nil {
			w.gameWidget.Update(gameStats)
			w.logger.LogGamePerformance(gameStats)
		}
	}
	
	// Collect Steam metrics
	if w.config.Widgets.ShowSteam {
		steamStats, err := w.steamCollector.Collect()
		if err == nil {
			w.steamWidget.Update(steamStats)
			w.logger.LogSteam(steamStats)
		}
	}
}

// ShowAndRun shows the window and runs the application
func (w *Window) ShowAndRun() {
	w.window.ShowAndRun()
}

// Close closes the window and cleans up resources
func (w *Window) Close() {
	if w.ticker != nil {
		w.ticker.Stop()
	}
	if w.logger != nil {
		w.logger.Close()
	}
	w.window.Close()
}

