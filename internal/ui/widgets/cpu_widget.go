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

// CPUWidget displays CPU usage metrics
type CPUWidget struct {
	widget.BaseWidget
	stats     *metrics.CPUStats
	theme     *theme.Theme
	title     *canvas.Text
	overall   *canvas.Text
	loadAvg   *canvas.Text
	coreBars  []*canvas.Rectangle
	coreLabels []*canvas.Text
	container *fyne.Container
}

// NewCPUWidget creates a new CPU widget
func NewCPUWidget(theme *theme.Theme) *CPUWidget {
	w := &CPUWidget{
		theme: theme,
		title: canvas.NewText("CPU", theme.TextColor),
		overall: canvas.NewText("Overall: 0%", theme.TextColor),
		loadAvg: canvas.NewText("Load: 0.00 0.00 0.00", theme.TextColor),
	}
	w.title.TextStyle = fyne.TextStyle{Bold: true}
	w.title.TextSize = 16
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer creates the renderer for the widget
func (w *CPUWidget) CreateRenderer() fyne.WidgetRenderer {
	w.container = container.NewVBox(
		w.title,
		w.overall,
		w.loadAvg,
	)

	// Add core bars
	for i := 0; i < 8; i++ { // Support up to 8 cores
		bar := canvas.NewRectangle(w.theme.BarColorLow)
		bar.SetMinSize(fyne.NewSize(200, 20))
		label := canvas.NewText(fmt.Sprintf("Core %d: 0%%", i), w.theme.TextColor)
		label.TextSize = 12

		coreContainer := container.NewBorder(nil, nil, label, nil, bar)
		w.coreBars = append(w.coreBars, bar)
		w.coreLabels = append(w.coreLabels, label)
		w.container.Add(coreContainer)
	}

	return &cpuWidgetRenderer{
		widget:    w,
		container: w.container,
	}
}

// Update updates the widget with new CPU stats
func (w *CPUWidget) Update(stats *metrics.CPUStats) {
	w.stats = stats
	if stats == nil {
		return
	}

	w.overall.Text = fmt.Sprintf("Overall: %.1f%%", stats.OverallPercent)
	w.overall.Refresh()

	w.loadAvg.Text = fmt.Sprintf("Load: %.2f %.2f %.2f", stats.LoadAvg1, stats.LoadAvg5, stats.LoadAvg15)
	w.loadAvg.Refresh()

	// Update core bars
	for i, percent := range stats.PerCorePercent {
		if i >= len(w.coreBars) {
			break
		}
		bar := w.coreBars[i]
		label := w.coreLabels[i]

		// Update bar color based on usage
		barColor := w.theme.GetBarColor(percent)
		bar.FillColor = barColor
		bar.Refresh()

		// Update label
		label.Text = fmt.Sprintf("Core %d: %.1f%%", i, percent)
		label.Refresh()

		// Resize bar to show percentage (simplified - would need custom layout)
		bar.SetMinSize(fyne.NewSize(float32(percent*2), 20))
	}
}

type cpuWidgetRenderer struct {
	widget    *CPUWidget
	container *fyne.Container
}

func (r *cpuWidgetRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *cpuWidgetRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *cpuWidgetRenderer) Refresh() {
	r.container.Refresh()
}

func (r *cpuWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *cpuWidgetRenderer) Destroy() {}

