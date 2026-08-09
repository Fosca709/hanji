package main

import (
	_ "embed"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/icon.svg
var iconData []byte

type hanjiEntryTheme struct {
	fyne.Theme
}

type hanjiTheme struct {
	fyne.Theme
}

var hanjiBackground = color.NRGBA{R: 0xF8, G: 0xE5, B: 0x8C, A: 0xFF}

func (t hanjiTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return hanjiBackground
	}
	if name == theme.ColorNamePrimary {
		return color.Black
	}
	return t.Theme.Color(name, variant)
}

func (t hanjiEntryTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameInputBackground {
		return hanjiBackground
	}
	if name == theme.ColorNameForeground {
		return color.Black
	}

	return t.Theme.Color(name, variant)
}

func (t hanjiEntryTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 16
	}
	return t.Theme.Size(name)
}

type hanjiEntry struct {
	widget.Entry
}

func newHanjiEntry() *container.ThemeOverride {
	entry := &hanjiEntry{}
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapWord
	entry.ExtendBaseWidget(entry)

	entryTheme := hanjiEntryTheme{Theme: theme.DefaultTheme()}
	return container.NewThemeOverride(entry, entryTheme)
}

func (e *hanjiEntry) CreateRenderer() fyne.WidgetRenderer {
	renderer := e.Entry.CreateRenderer()
	e.ExtendBaseWidget(e)

	hanjiRenderer := &hanjiEntryRenderer{WidgetRenderer: renderer}
	hanjiRenderer.hideBorder()
	return hanjiRenderer
}

type hanjiEntryRenderer struct {
	fyne.WidgetRenderer
}

func (r *hanjiEntryRenderer) Layout(size fyne.Size) {
	r.WidgetRenderer.Layout(size)
	r.hideBorder()
}

func (r *hanjiEntryRenderer) Refresh() {
	r.WidgetRenderer.Refresh()
	r.hideBorder()
}

func (r *hanjiEntryRenderer) hideBorder() {
	objects := r.WidgetRenderer.Objects()
	if len(objects) < 2 {
		return
	}

	if border, ok := objects[1].(*canvas.Rectangle); ok {
		border.StrokeWidth = 0
	}
}

func main() {
	// Fyne's GLFW backend needs XIM configured to receive composed Korean input
	// from fcitx5. This is intentionally fcitx5-specific for this personal app.
	if os.Getenv("XMODIFIERS") == "" {
		_ = os.Setenv("XMODIFIERS", "@im=fcitx")
	}

	a := app.NewWithID("Hanji")
	a.SetIcon(fyne.NewStaticResource("icon.svg", iconData))
	a.Settings().SetTheme(hanjiTheme{Theme: theme.DefaultTheme()})
	w := a.NewWindow("Hanji")

	note := newHanjiEntry()

	w.SetContent(note)
	w.Resize(fyne.NewSize(350, 350))
	w.ShowAndRun()
}
