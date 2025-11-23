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

// GameWidget displays game performance metrics
type GameWidget struct {
	widget.BaseWidget
	stats     *metrics.GamePerformanceStats
	theme     *theme.Theme
	title     *canvas.Text
	gameName  *canvas.Text
	fpsText   *canvas.Text
	frameTimeText *canvas.Text
	container *fyne.Container
}

// NewGameWidget creates a new game widget
func NewGameWidget(theme *theme.Theme) *GameWidget {
	w := &GameWidget{
		theme: theme,
		title: canvas.NewText("Game Performance", theme.TextColor),
		gameName: canvas.NewText("Game: None", theme.TextColor),
		fpsText: canvas.NewText("FPS: 0", theme.TextColor),
		frameTimeText: canvas.NewText("Frame Time: 0.00 ms", theme.TextColor),
	}
	w.title.TextStyle = fyne.TextStyle{Bold: true}
	w.title.TextSize = 16
	w.gameName.TextSize = 14
	w.fpsText.TextSize = 14
	w.frameTimeText.TextSize = 12
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer creates the renderer for the widget
func (w *GameWidget) CreateRenderer() fyne.WidgetRenderer {
	w.container = container.NewVBox(
		w.title,
		w.gameName,
		w.fpsText,
		w.frameTimeText,
	)

	return &gameWidgetRenderer{
		widget:    w,
		container: w.container,
	}
}

// Update updates the widget with new game performance stats
func (w *GameWidget) Update(stats *metrics.GamePerformanceStats) {
	w.stats = stats
	if stats == nil {
		return
	}

	if stats.GameName != "" {
		w.gameName.Text = fmt.Sprintf("Game: %s", stats.GameName)
	} else {
		w.gameName.Text = "Game: None"
	}
	w.gameName.Refresh()

	w.fpsText.Text = fmt.Sprintf("FPS: %.1f", stats.FPS)
	w.fpsText.Refresh()

	w.frameTimeText.Text = fmt.Sprintf("Frame Time: %.2f ms (Min: %.2f ms, Max: %.2f ms)",
		stats.FrameTime, stats.FrameTimeMin, stats.FrameTimeMax)
	w.frameTimeText.Refresh()
}

type gameWidgetRenderer struct {
	widget    *GameWidget
	container *fyne.Container
}

func (r *gameWidgetRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *gameWidgetRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *gameWidgetRenderer) Refresh() {
	r.container.Refresh()
}

func (r *gameWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *gameWidgetRenderer) Destroy() {}

