package main

import (
	_ "embed"
	"image/color"
	"log"
	"os"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
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
	if name == theme.ColorNameMenuBackground {
		return color.White
	}
	if name == theme.ColorNameForeground {
		return color.Black
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
	alwaysOnTop bool
}

type menuToggleIndicator bool

func (i menuToggleIndicator) ShortcutName() string {
	return "AlwaysOnTopIndicator"
}

func (i menuToggleIndicator) Key() fyne.KeyName {
	if i {
		return fyne.KeyName("▣")
	}
	return fyne.KeyName("□")
}

func (menuToggleIndicator) Mod() fyne.KeyModifier {
	return 0
}

func (e *hanjiEntry) toggleAlwaysOnTop() {
	enabled := !e.alwaysOnTop
	if err := setAlwaysOnTop(enabled); err != nil {
		log.Printf("could not change always-on-top state: %v", err)
		return
	}
	e.alwaysOnTop = enabled
}

func (e *hanjiEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if keyboard, ok := shortcut.(fyne.KeyboardShortcut); ok &&
		keyboard.Key() == fyne.KeyQ && keyboard.Mod() == fyne.KeyModifierControl {
		e.toggleAlwaysOnTop()
		return
	}
	e.Entry.TypedShortcut(shortcut)
}

func (e *hanjiEntry) TappedSecondary(pe *fyne.PointEvent) {
	canUndo, canRedo := e.undoAvailability()

	alwaysOnTopItem := fyne.NewMenuItem("Always on Top", e.toggleAlwaysOnTop)
	alwaysOnTopItem.Shortcut = menuToggleIndicator(e.alwaysOnTop)
	undoItem := fyne.NewMenuItem("Undo", e.Undo)
	undoItem.Disabled = !canUndo
	redoItem := fyne.NewMenuItem("Redo", e.Redo)
	redoItem.Disabled = !canRedo

	app := fyne.CurrentApp()
	canvas := app.Driver().CanvasForObject(e)
	canvas.Focus(e)
	widget.ShowPopUpMenuAtPosition(
		fyne.NewMenu("", alwaysOnTopItem, undoItem, redoItem),
		canvas,
		pe.AbsolutePosition,
	)
}

func (e *hanjiEntry) undoAvailability() (canUndo, canRedo bool) {
	stack := reflect.ValueOf(&e.Entry).Elem().FieldByName("undoStack")
	if !stack.IsValid() {
		return false, false
	}

	index := stack.FieldByName("index")
	items := stack.FieldByName("items")
	if !index.IsValid() || !items.IsValid() {
		return false, false
	}

	position := int(index.Int())
	return position > 0, position < items.Len()
}

func newHanjiEntry() (*container.ThemeOverride, *hanjiEntry) {
	entry := &hanjiEntry{}
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapBreak
	entry.ExtendBaseWidget(entry)

	entryTheme := hanjiEntryTheme{Theme: theme.DefaultTheme()}
	return container.NewThemeOverride(entry, entryTheme), entry
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

	note, entry := newHanjiEntry()
	alwaysOnTopShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}
	w.Canvas().AddShortcut(alwaysOnTopShortcut, func(fyne.Shortcut) {
		entry.toggleAlwaysOnTop()
	})

	w.SetContent(note)
	w.Resize(fyne.NewSize(350, 350))
	w.ShowAndRun()
}
