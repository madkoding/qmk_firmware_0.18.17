package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/qmk-tui/build"
	"github.com/qmk-tui/core"
)

type App struct {
	View        *tview.Application
	Pages       *tview.Pages
	KBList      *KeyboardListView
	Canvas      *KeyboardCanvas
	ActionPanel *ActionPanel
	LogView     *tview.TextView
	Compiler    *build.Compiler
	KeycodePicker *KeycodePicker
	CurrentKB   *core.Keyboard
	CurrentKM   *core.Keymap
	CurrentLayer int
}

func Run() error {
	app := &App{
		View:     tview.NewApplication(),
		Compiler: build.NewCompiler(core.GetQMKPath()),
	}

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorBlack,
		ContrastBackgroundColor:     tcell.ColorBlack,
		MoreContrastBackgroundColor: tcell.ColorBlack,
		BorderColor:                 tcell.NewHexColor(0xf9e2af), // yellow
		TitleColor:                  tcell.NewHexColor(0x89b4fa), // blue
		GraphicsColor:               tcell.NewHexColor(0x6c7086), // overlay0
		PrimaryTextColor:            tcell.NewHexColor(0xcdd6f4), // text
		SecondaryTextColor:          tcell.NewHexColor(0xbac2de), // subtext1
		TertiaryTextColor:           tcell.NewHexColor(0x6c7086), // overlay0
		InverseTextColor:            tcell.NewHexColor(0x1e1e2e),
		ContrastSecondaryTextColor:  tcell.NewHexColor(0x94e2d5), // teal
	}

	if err := app.setup(); err != nil {
		return err
	}

	return app.View.SetRoot(app.Pages, true).Run()
}

func (app *App) setup() error {
	app.KBList = NewKeyboardListView(func(kb *core.Keyboard, km string) {
		app.SelectKeyboard(kb, km)
	})

	app.Canvas = NewKeyboardCanvas(func(key *core.Key) {
		app.ShowKeycodePickerForKey(key)
	})

	app.ActionPanel = NewActionPanel()
	app.LogView = app.ActionPanel.LogView

	app.KeycodePicker = NewKeycodePicker(func(kc string) {
		app.Canvas.SetKeycode(kc)
		app.Pages.HidePage("keypicker")
	})

	app.ActionPanel.BuildBtn.SetSelectedFunc(func() {
		app.Build()
	})

	app.ActionPanel.FlashBtn.SetSelectedFunc(func() {
		app.Flash()
	})

	app.ActionPanel.CopyBtn.SetSelectedFunc(func() {
		app.ShowCopyDialog()
	})

	app.ActionPanel.SaveBtn.SetSelectedFunc(func() {
		app.Save()
	})

	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	leftPanel.AddItem(app.KBList, 0, 1, true)

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	rightPanel.AddItem(app.Canvas, 0, 2, false)
	rightPanel.AddItem(app.ActionPanel, 0, 1, false)

	mainContent := tview.NewFlex().
		AddItem(leftPanel, 25, 0, true).
		AddItem(rightPanel, 0, 1, false)

	app.Pages = tview.NewPages()
	app.Pages.AddPage("main", mainContent, true, true)

	app.Pages.SetInputCapture(app.handleInput)

	app.View.SetFocus(app.KBList)

	return nil
}

func (app *App) cycleFocus() {
	focusOrder := []tview.Primitive{app.KBList, app.Canvas, app.ActionPanel.BuildBtn}
	current := app.View.GetFocus()
	for i, p := range focusOrder {
		if current == p {
			app.View.SetFocus(focusOrder[(i+1)%len(focusOrder)])
			return
		}
	}
	app.View.SetFocus(app.KBList)
}

func (app *App) isCanvasFocused() bool {
	f := app.View.GetFocus()
	if f == app.Canvas || f == app.Canvas.TextView {
		return true
	}
	_, ok := f.(*KeyboardCanvas)
	return ok
}

func (app *App) isKBListFocused() bool {
	f := app.View.GetFocus()
	if f == app.KBList || f == app.KBList.TreeView {
		return true
	}
	_, ok := f.(*KeyboardListView)
	return ok
}

func (app *App) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		app.View.Stop()
		return nil
	case tcell.KeyTab:
		app.cycleFocus()
		return nil
	case tcell.KeyEnter:
		if app.isCanvasFocused() {
			sel := app.Canvas.GetSelectedKey()
			if sel != nil {
				app.ShowKeycodePickerForKey(sel)
				return nil
			}
		}
		if app.isKBListFocused() {
			app.KBList.SelectCurrentNode()
			return nil
		}
		return event
	case tcell.KeyUp:
		if app.isCanvasFocused() {
			app.Canvas.Navigate("up")
			return nil
		}
		return event
	case tcell.KeyDown:
		if app.isCanvasFocused() {
			app.Canvas.Navigate("down")
			return nil
		}
		return event
	case tcell.KeyLeft:
		if app.isCanvasFocused() {
			app.Canvas.Navigate("left")
			return nil
		}
		return event
	case tcell.KeyRight:
		if app.isCanvasFocused() {
			app.Canvas.Navigate("right")
			return nil
		}
		return event
	case tcell.KeyRune:
		switch event.Rune() {
		case 's', 'S':
			if app.CurrentKM != nil {
				app.Save()
			}
		case 'b', 'B':
			app.Build()
		case 'f', 'F':
			app.Flash()
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			app.SwitchLayer(int(event.Rune() - '1'))
		}
	}

	return event
}

func (app *App) SwitchLayer(layer int) {
	if app.CurrentKM == nil || layer >= len(app.CurrentKM.Layers) {
		return
	}
	app.CurrentLayer = layer
	app.Canvas.SetLayer(layer)
	app.Log(fmt.Sprintf("Switched to layer %d (%s)", layer, app.CurrentKM.Layers[layer].Name))
}

func (app *App) Log(msg string) {
	app.LogView.SetText(app.LogView.GetText(true) + msg + "\n")
}

func (app *App) SelectKeyboard(kb *core.Keyboard, km string) {
	app.CurrentKB = kb
	app.Log(fmt.Sprintf("Selected: %s:%s", kb.Name, km))

	kmPtr, err := core.LoadKeymap(kb, km)
	if err != nil {
		app.Log(fmt.Sprintf("Error loading keymap: %v", err))
		return
	}

	app.CurrentKM = kmPtr
	app.CurrentLayer = 0
	app.Canvas.SetKeyboard(kb, kmPtr)
 	app.Canvas.SetLayer(0)

	app.Log(fmt.Sprintf("Loaded %d layers", len(kmPtr.Layers)))
	for i, layer := range kmPtr.Layers {
		app.Log(fmt.Sprintf("Layer %d: %s", i, layer.Name))
	}
}

func (app *App) Build() {
	if app.CurrentKB == nil || app.CurrentKM == nil {
		app.Log("No keyboard/keymap selected")
		return
	}

	app.Log(fmt.Sprintf("Building %s:%s...", app.CurrentKB.Name, app.CurrentKM.KeymapName))

	output, err := app.Compiler.Build(app.CurrentKB.Name, app.CurrentKM.KeymapName)
	if err != nil {
		app.Log(fmt.Sprintf("Build failed: %v", err))
		app.Log(output)
		return
	}

	if output == "" {
		app.Log("Build successful!")
	} else {
		app.Log(output)
	}
}

func (app *App) Flash() {
	if app.CurrentKB == nil || app.CurrentKM == nil {
		app.Log("No keyboard/keymap selected")
		return
	}

	app.Log(fmt.Sprintf("Flashing %s:%s...", app.CurrentKB.Name, app.CurrentKM.KeymapName))

	output, err := app.Compiler.Flash(app.CurrentKB.Name, app.CurrentKM.KeymapName)
	if err != nil {
		app.Log(fmt.Sprintf("Flash failed: %v", err))
		app.Log(output)
		return
	}

	app.Log("Flash successful!")
	if output != "" {
		app.Log(output)
	}
}

func (app *App) ShowCopyDialog() {
	if app.CurrentKB == nil {
		app.Log("No keyboard selected")
		return
	}

	header := tview.NewTextView().SetText("Select source keymap to copy from:")

	list := tview.NewList()
	for _, km := range app.CurrentKB.Keymaps {
		if km != app.CurrentKM.KeymapName {
			list.AddItem(km, "", 0, nil)
		}
	}

	if list.GetItemCount() == 0 {
		app.Log("No other keymaps available to copy from")
		return
	}

	list.AddItem("Cancel", "", 0, nil)

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "Cancel" {
			app.Pages.HidePage("copy_dialog")
			return
		}
		app.CopyFromKeymap(mainText)
		app.Pages.HidePage("copy_dialog")
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(header, 1, 0, false)
	flex.AddItem(list, 10, 1, true)

	app.Pages.AddPage("copy_dialog", flex, true, false)
}

func (app *App) CopyFromKeymap(srcKeymapName string) {
	if app.CurrentKB == nil || app.CurrentKM == nil {
		return
	}

	srcKM, err := core.LoadKeymap(app.CurrentKB, srcKeymapName)
	if err != nil {
		app.Log(fmt.Sprintf("Error loading source keymap: %v", err))
		return
	}

	newKM, err := core.CopyKeymap(srcKM, app.CurrentKB, app.CurrentKM.KeymapName)
	if err != nil {
		app.Log(fmt.Sprintf("Error copying keymap: %v", err))
		return
	}

	app.CurrentKM = newKM
	app.Canvas.SetKeyboard(app.CurrentKB, newKM)
	app.Log(fmt.Sprintf("Copied layout from %s", srcKeymapName))
}

func (app *App) Save() {
	if app.CurrentKM == nil {
		app.Log("No keymap loaded")
		return
	}

	if err := app.CurrentKM.Save(); err != nil {
		app.Log(fmt.Sprintf("Save failed: %v", err))
		return
	}

	app.Log("Saved successfully!")
}

func (app *App) ShowKeycodePickerForKey(key *core.Key) {
	if app.CurrentKM == nil || key == nil {
		return
	}

	kc := key.Keycode
	if kc == "" {
		kc = "KC_NO"
	}

	app.KeycodePicker.ShowForKey(key, kc)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("Editing key [row:%d, col:%d] current:[%s]", key.Row, key.Col, kc)), 3, 0, false)
	flex.AddItem(app.KeycodePicker.GetSearchBar(), 2, 0, true)
	flex.AddItem(app.KeycodePicker.GetList(), 12, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			selectedKC := app.KeycodePicker.GetSelected()
			_ = key
			app.Canvas.SetKeycode(selectedKC)
			app.Pages.HidePage("keypicker")
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			app.Pages.HidePage("keypicker")
			return nil
		}
		return event
	})

	app.Pages.AddPage("keypicker", flex, true, false)
}
