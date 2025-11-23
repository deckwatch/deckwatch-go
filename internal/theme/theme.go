package theme

import (
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Theme provides btop-inspired color scheme
type Theme struct {
	BackgroundColor color.Color
	TextColor       color.Color
	BarColor        color.Color
	BarColorHigh    color.Color
	BarColorMedium  color.Color
	BarColorLow     color.Color
}

// NewTheme creates a new theme from hex color strings
func NewTheme(bgColor, textColor, barColor, barHigh, barMed, barLow string) *Theme {
	return &Theme{
		BackgroundColor: parseColor(bgColor),
		TextColor:       parseColor(textColor),
		BarColor:        parseColor(barColor),
		BarColorHigh:    parseColor(barHigh),
		BarColorMedium: parseColor(barMed),
		BarColorLow:     parseColor(barLow),
	}
}

// parseColor parses a hex color string to color.Color
func parseColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}

// GetBarColor returns the appropriate bar color based on percentage
func (t *Theme) GetBarColor(percent float64) color.Color {
	if percent >= 80 {
		return t.BarColorHigh
	} else if percent >= 50 {
		return t.BarColorMedium
	} else {
		return t.BarColorLow
	}
}

// DefaultTheme returns the default btop-inspired theme
func DefaultTheme() *Theme {
	return NewTheme(
		"#1e1e2e", // Background
		"#cdd6f4", // Text
		"#89b4fa", // Bar
		"#f38ba8", // High (red)
		"#fab387", // Medium (orange)
		"#a6e3a1", // Low (green)
	)
}

// ApplyTheme applies the theme to a Fyne app
func ApplyTheme(app fyne.App, t *Theme) {
	// Create a custom theme that uses our colors
	customTheme := &customTheme{
		Theme: theme.DefaultTheme(),
		colors: t,
	}
	app.Settings().SetTheme(customTheme)
}

type customTheme struct {
	fyne.Theme
	colors *Theme
}

func (c *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return c.colors.BackgroundColor
	case theme.ColorNameForeground:
		return c.colors.TextColor
	default:
		return c.Theme.Color(name, variant)
	}
}

