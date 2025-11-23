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

// SteamWidget displays Steam-specific metrics
type SteamWidget struct {
	widget.BaseWidget
	stats     *metrics.SteamStats
	theme     *theme.Theme
	title     *canvas.Text
	downloadText *canvas.Text
	uploadText   *canvas.Text
	libraryText  *canvas.Text
	downloadsText *canvas.Text
	container *fyne.Container
}

// NewSteamWidget creates a new Steam widget
func NewSteamWidget(theme *theme.Theme) *SteamWidget {
	w := &SteamWidget{
		theme: theme,
		title: canvas.NewText("Steam", theme.TextColor),
		downloadText: canvas.NewText("Download: 0 MB/s", theme.TextColor),
		uploadText: canvas.NewText("Upload: 0 MB/s", theme.TextColor),
		libraryText: canvas.NewText("Library: 0 games, 0 GB", theme.TextColor),
		downloadsText: canvas.NewText("Active Downloads: 0", theme.TextColor),
	}
	w.title.TextStyle = fyne.TextStyle{Bold: true}
	w.title.TextSize = 16
	w.downloadText.TextSize = 14
	w.uploadText.TextSize = 14
	w.libraryText.TextSize = 12
	w.downloadsText.TextSize = 12
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer creates the renderer for the widget
func (w *SteamWidget) CreateRenderer() fyne.WidgetRenderer {
	w.container = container.NewVBox(
		w.title,
		w.downloadText,
		w.uploadText,
		w.libraryText,
		w.downloadsText,
	)

	return &steamWidgetRenderer{
		widget:    w,
		container: w.container,
	}
}

// Update updates the widget with new Steam stats
func (w *SteamWidget) Update(stats *metrics.SteamStats) {
	w.stats = stats
	if stats == nil {
		return
	}

	downloadMB := stats.DownloadSpeed / (1024 * 1024)
	uploadMB := stats.UploadSpeed / (1024 * 1024)
	libraryGB := float64(stats.LibrarySize) / (1024 * 1024 * 1024)

	w.downloadText.Text = fmt.Sprintf("Download: %.2f MB/s", downloadMB)
	w.downloadText.Refresh()

	w.uploadText.Text = fmt.Sprintf("Upload: %.2f MB/s", uploadMB)
	w.uploadText.Refresh()

	w.libraryText.Text = fmt.Sprintf("Library: %d games, %.2f GB", stats.InstalledGames, libraryGB)
	w.libraryText.Refresh()

	w.downloadsText.Text = fmt.Sprintf("Active Downloads: %d", stats.ActiveDownloads)
	w.downloadsText.Refresh()

	// Add download progress if available
	if len(stats.DownloadProgress) > 0 {
		// Could add progress bars here for each download
	}
}

type steamWidgetRenderer struct {
	widget    *SteamWidget
	container *fyne.Container
}

func (r *steamWidgetRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *steamWidgetRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *steamWidgetRenderer) Refresh() {
	r.container.Refresh()
}

func (r *steamWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *steamWidgetRenderer) Destroy() {}

