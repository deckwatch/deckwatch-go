package widgets

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/steam-os-monitor/monitor/internal/theme"
	"github.com/steam-os-monitor/monitor/pkg/metrics"
)

// MemoryWidget displays memory usage metrics
type MemoryWidget struct {
	widget.BaseWidget
	stats     *metrics.MemoryStats
	theme     *theme.Theme
	title     *canvas.Text
	ramText   *canvas.Text
	swapText  *canvas.Text
	ramBar    *canvas.Rectangle
	swapBar   *canvas.Rectangle
	container *fyne.Container
}

// NewMemoryWidget creates a new memory widget
func NewMemoryWidget(theme *theme.Theme) *MemoryWidget {
	w := &MemoryWidget{
		theme: theme,
		title: canvas.NewText("Memory", theme.TextColor),
		ramText: canvas.NewText("RAM: 0 / 0 GB (0%)", theme.TextColor),
		swapText: canvas.NewText("Swap: 0 / 0 GB (0%)", theme.TextColor),
		ramBar: canvas.NewRectangle(theme.BarColorLow),
		swapBar: canvas.NewRectangle(theme.BarColorLow),
	}
	w.title.TextStyle = fyne.TextStyle{Bold: true}
	w.title.TextSize = 16
	w.ramBar.SetMinSize(fyne.NewSize(300, 30))
	w.swapBar.SetMinSize(fyne.NewSize(300, 20))
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer creates the renderer for the widget
func (w *MemoryWidget) CreateRenderer() fyne.WidgetRenderer {
	w.container = container.NewVBox(
		w.title,
		w.ramText,
		container.NewWithoutLayout(w.ramBar),
		w.swapText,
		container.NewWithoutLayout(w.swapBar),
	)

	return &memoryWidgetRenderer{
		widget:    w,
		container: w.container,
	}
}

// Update updates the widget with new memory stats
func (w *MemoryWidget) Update(stats *metrics.MemoryStats) {
	w.stats = stats
	if stats == nil {
		return
	}

	// Update RAM text
	ramUsedGB := float64(stats.Used) / (1024 * 1024 * 1024)
	ramTotalGB := float64(stats.Total) / (1024 * 1024 * 1024)
	w.ramText.Text = fmt.Sprintf("RAM: %.2f / %.2f GB (%.1f%%)", ramUsedGB, ramTotalGB, stats.UsedPercent)
	w.ramText.Refresh()

	// Update RAM bar color and size
	ramColor := w.theme.GetBarColor(stats.UsedPercent)
	w.ramBar.FillColor = ramColor
	w.ramBar.SetMinSize(fyne.NewSize(float32(stats.UsedPercent*3), 30))
	w.ramBar.Refresh()

	// Update Swap text
	if stats.SwapTotal > 0 {
		swapUsedGB := float64(stats.SwapUsed) / (1024 * 1024 * 1024)
		swapTotalGB := float64(stats.SwapTotal) / (1024 * 1024 * 1024)
		w.swapText.Text = fmt.Sprintf("Swap: %.2f / %.2f GB (%.1f%%)", swapUsedGB, swapTotalGB, stats.SwapPercent)
		w.swapText.Refresh()

		// Update Swap bar
		swapColor := w.theme.GetBarColor(stats.SwapPercent)
		w.swapBar.FillColor = swapColor
		w.swapBar.SetMinSize(fyne.NewSize(float32(stats.SwapPercent*3), 20))
		w.swapBar.Refresh()
	} else {
		w.swapText.Text = "Swap: Not Available"
		w.swapText.Refresh()
	}
}

type memoryWidgetRenderer struct {
	widget    *MemoryWidget
	container *fyne.Container
}

func (r *memoryWidgetRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *memoryWidgetRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *memoryWidgetRenderer) Refresh() {
	r.container.Refresh()
}

func (r *memoryWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *memoryWidgetRenderer) Destroy() {}

