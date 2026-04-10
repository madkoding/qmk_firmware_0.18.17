package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ActionPanel struct {
	*tview.Flex
	BuildBtn  *tview.Button
	FlashBtn  *tview.Button
	CopyBtn   *tview.Button
	SaveBtn   *tview.Button
	LogView   *tview.TextView
	HelpView  *tview.TextView
}

func NewActionPanel() *ActionPanel {
	panel := &ActionPanel{
		Flex:     tview.NewFlex(),
		BuildBtn: tview.NewButton("[Build]"),
		FlashBtn: tview.NewButton("[Flash]"),
		CopyBtn:  tview.NewButton("[Copy]"),
		SaveBtn:  tview.NewButton("[Save]"),
		LogView:  tview.NewTextView(),
		HelpView: tview.NewTextView(),
	}

	panel.SetTitle(" Actions ")
	panel.SetBorder(true)
	panel.SetBorderColor(tcell.NewHexColor(0xf9e2af))
	panel.SetDirection(tview.FlexRow)

	btnRow := tview.NewFlex()
	btnRow.AddItem(panel.BuildBtn, 0, 1, false)
	btnRow.AddItem(panel.FlashBtn, 0, 1, false)
	btnRow.AddItem(panel.CopyBtn, 0, 1, false)
	btnRow.AddItem(panel.SaveBtn, 0, 1, false)

	panel.HelpView.SetDynamicColors(true)
	panel.HelpView.SetText(`[#89b4fa]Controls:[white]
  [#94e2d5]Tab[white]       Change panel
  [#94e2d5]Arrows[white]    Navigate/Select
  [#94e2d5]Enter[white]     Edit key
  [#94e2d5]1-4[white]       Switch layer
  [#a6e3a1]s[white]         Save
  [#f9e2af]b[white]         Build
  [#fab387]f[white]         Flash
  [#f38ba8]Ctrl+C[white]    Exit`)
	panel.HelpView.SetBorder(true)
	panel.HelpView.SetTitle(" Help ")
	panel.HelpView.SetBorderColor(tcell.NewHexColor(0x94e2d5))
	panel.HelpView.SetBackgroundColor(tcell.ColorBlack)

	panel.LogView.SetTitle(" Log ")
	panel.LogView.SetBorder(true)
	panel.LogView.SetDynamicColors(true)
	panel.LogView.SetScrollable(true)
	panel.LogView.SetBackgroundColor(tcell.ColorBlack)

	panel.AddItem(panel.HelpView, 11, 0, false)
	panel.AddItem(btnRow, 3, 0, false)
	panel.AddItem(panel.LogView, 0, 1, false)

	return panel
}

func (p *ActionPanel) Log(msg string) {
	p.LogView.SetText(p.LogView.GetText(true) + msg + "\n")
}

func (p *ActionPanel) ClearLog() {
	p.LogView.SetText("")
}
