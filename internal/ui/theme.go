package ui

// Re-export theme types for convenience
import (
	"fyne.io/fyne/v2"
	"github.com/steam-os-monitor/monitor/internal/theme"
)

// Theme is re-exported from theme package
type Theme = theme.Theme

// DefaultTheme returns the default theme
func DefaultTheme() *Theme {
	return theme.DefaultTheme()
}

// ApplyTheme applies the theme to a Fyne app
func ApplyTheme(app fyne.App, t *Theme) {
	theme.ApplyTheme(app, t)
}
