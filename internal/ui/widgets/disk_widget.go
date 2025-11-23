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

// DiskWidget displays disk usage metrics
type DiskWidget struct {
	widget.BaseWidget
	stats     []*metrics.DiskStats
	theme     *theme.Theme
	title     *canvas.Text
	disks     []*diskEntry
	container *fyne.Container
}

type diskEntry struct {
	label    *canvas.Text
	bar      *canvas.Rectangle
	ioText   *canvas.Text
	container *fyne.Container
}

// NewDiskWidget creates a new disk widget
func NewDiskWidget(theme *theme.Theme) *DiskWidget {
	w := &DiskWidget{
		theme: theme,
		title: canvas.NewText("Disk", theme.TextColor),
		disks: make([]*diskEntry, 0),
	}
	w.title.TextStyle = fyne.TextStyle{Bold: true}
	w.title.TextSize = 16
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer creates the renderer for the widget
func (w *DiskWidget) CreateRenderer() fyne.WidgetRenderer {
	w.container = container.NewVBox(w.title)
	return &diskWidgetRenderer{
		widget:    w,
		container: w.container,
	}
}

// Update updates the widget with new disk stats
func (w *DiskWidget) Update(stats []*metrics.DiskStats) {
	w.stats = stats
	if stats == nil {
		return
	}

	// Ensure we have enough disk entries
	for len(w.disks) < len(stats) {
		entry := &diskEntry{
			label: canvas.NewText("", w.theme.TextColor),
			bar: canvas.NewRectangle(w.theme.BarColorLow),
			ioText: canvas.NewText("", w.theme.TextColor),
		}
		entry.label.TextSize = 12
		entry.ioText.TextSize = 10
		entry.bar.SetMinSize(fyne.NewSize(300, 20))
		entry.container = container.NewVBox(
			entry.label,
			container.NewWithoutLayout(entry.bar),
			entry.ioText,
		)
		w.disks = append(w.disks, entry)
		w.container.Add(entry.container)
	}

	// Update each disk entry
	for i, stat := range stats {
		if i >= len(w.disks) {
			break
		}
		entry := w.disks[i]

		// Format size
		usedGB := float64(stat.Used) / (1024 * 1024 * 1024)
		totalGB := float64(stat.Total) / (1024 * 1024 * 1024)
		entry.label.Text = fmt.Sprintf("%s (%s): %.2f / %.2f GB (%.1f%%)",
			stat.Device, stat.MountPoint, usedGB, totalGB, stat.UsedPercent)
		entry.label.Refresh()

		// Update bar
		barColor := w.theme.GetBarColor(stat.UsedPercent)
		entry.bar.FillColor = barColor
		entry.bar.SetMinSize(fyne.NewSize(float32(stat.UsedPercent*3), 20))
		entry.bar.Refresh()

		// Update IO stats
		readMB := float64(stat.ReadBytes) / (1024 * 1024)
		writeMB := float64(stat.WriteBytes) / (1024 * 1024)
		entry.ioText.Text = fmt.Sprintf("  Read: %.2f MB | Write: %.2f MB | IOPS: R:%d W:%d",
			readMB, writeMB, stat.ReadIOPS, stat.WriteIOPS)
		entry.ioText.Refresh()
	}
}

type diskWidgetRenderer struct {
	widget    *DiskWidget
	container *fyne.Container
}

func (r *diskWidgetRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *diskWidgetRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *diskWidgetRenderer) Refresh() {
	r.container.Refresh()
}

func (r *diskWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *diskWidgetRenderer) Destroy() {}

