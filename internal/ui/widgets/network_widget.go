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

// NetworkWidget displays network statistics
type NetworkWidget struct {
	widget.BaseWidget
	stats     []*metrics.NetworkStats
	theme     *theme.Theme
	title     *canvas.Text
	interfaces []*networkEntry
	container *fyne.Container
}

type networkEntry struct {
	label     *canvas.Text
	sentText  *canvas.Text
	recvText  *canvas.Text
	container *fyne.Container
}

// NewNetworkWidget creates a new network widget
func NewNetworkWidget(theme *theme.Theme) *NetworkWidget {
	w := &NetworkWidget{
		theme: theme,
		title: canvas.NewText("Network", theme.TextColor),
		interfaces: make([]*networkEntry, 0),
	}
	w.title.TextStyle = fyne.TextStyle{Bold: true}
	w.title.TextSize = 16
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer creates the renderer for the widget
func (w *NetworkWidget) CreateRenderer() fyne.WidgetRenderer {
	w.container = container.NewVBox(w.title)
	return &networkWidgetRenderer{
		widget:    w,
		container: w.container,
	}
}

// Update updates the widget with new network stats
func (w *NetworkWidget) Update(stats []*metrics.NetworkStats) {
	w.stats = stats
	if stats == nil {
		return
	}

	// Ensure we have enough interface entries
	for len(w.interfaces) < len(stats) {
		entry := &networkEntry{
			label: canvas.NewText("", w.theme.TextColor),
			sentText: canvas.NewText("", w.theme.TextColor),
			recvText: canvas.NewText("", w.theme.TextColor),
		}
		entry.label.TextStyle = fyne.TextStyle{Bold: true}
		entry.label.TextSize = 12
		entry.sentText.TextSize = 11
		entry.recvText.TextSize = 11
		entry.container = container.NewVBox(
			entry.label,
			entry.sentText,
			entry.recvText,
		)
		w.interfaces = append(w.interfaces, entry)
		w.container.Add(entry.container)
	}

	// Update each interface
	for i, stat := range stats {
		if i >= len(w.interfaces) {
			break
		}
		entry := w.interfaces[i]

		entry.label.Text = fmt.Sprintf("Interface: %s", stat.Interface)
		entry.label.Refresh()

		// Format bytes
		sentMB := float64(stat.BytesSent) / (1024 * 1024)
		recvMB := float64(stat.BytesRecv) / (1024 * 1024)
		sentSpeedMB := stat.SpeedSent / (1024 * 1024)
		recvSpeedMB := stat.SpeedRecv / (1024 * 1024)

		entry.sentText.Text = fmt.Sprintf("  Sent: %.2f MB (%.2f MB/s) | Packets: %d",
			sentMB, sentSpeedMB, stat.PacketsSent)
		entry.sentText.Refresh()

		entry.recvText.Text = fmt.Sprintf("  Recv: %.2f MB (%.2f MB/s) | Packets: %d",
			recvMB, recvSpeedMB, stat.PacketsRecv)
		entry.recvText.Refresh()
	}
}

type networkWidgetRenderer struct {
	widget    *NetworkWidget
	container *fyne.Container
}

func (r *networkWidgetRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *networkWidgetRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *networkWidgetRenderer) Refresh() {
	r.container.Refresh()
}

func (r *networkWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *networkWidgetRenderer) Destroy() {}

