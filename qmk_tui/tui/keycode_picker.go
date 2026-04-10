package tui

import (
	"github.com/rivo/tview"
	"github.com/qmk-tui/core"
)

type KeycodePicker struct {
	Modal      *tview.Modal
	onSelect   func(kc string)
	searchBar  *tview.InputField
	list       *tview.List
	categories []core.KeycodeCategory
	selectedKC string
	key        *core.Key
}

func NewKeycodePicker(onSelect func(kc string)) *KeycodePicker {
	picker := &KeycodePicker{
		onSelect:   onSelect,
		categories: core.KEYCODE_CATEGORIES,
		selectedKC: "KC_NO",
	}

	picker.searchBar = tview.NewInputField()
	picker.searchBar.SetLabel("Search: ")
	picker.searchBar.SetPlaceholder("Type to filter...")
	picker.searchBar.SetChangedFunc(func(text string) {
		picker.populateList(text)
	})

	picker.list = tview.NewList()
	picker.list.ShowSecondaryText(false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(picker.searchBar, 2, 0, true)
	flex.AddItem(picker.list, 12, 1, false)

	picker.populateList("")

	picker.list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		picker.selectedKC = mainText
	})

	return picker
}

func (p *KeycodePicker) populateList(filter string) {
	p.list.Clear()

	for _, cat := range p.categories {
		for _, kc := range cat.Keycodes {
			if filter == "" || containsIgnoreCase(kc, filter) {
				p.list.AddItem(kc, "", 0, nil)
			}
		}
	}

	if p.list.GetItemCount() > 0 {
		p.list.SetCurrentItem(0)
		mainText, _ := p.list.GetItemText(0)
		p.selectedKC = mainText
	}
}

func (p *KeycodePicker) ShowForKey(key *core.Key, currentKC string) {
	p.key = key
	p.selectedKC = currentKC
	if p.selectedKC == "" {
		p.selectedKC = "KC_NO"
	}
	p.searchBar.SetText("")
	p.populateList("")
}

func (p *KeycodePicker) GetSelected() string {
	idx := p.list.GetCurrentItem()
	if idx >= 0 && idx < len(core.ALL_KEYCODES) {
		mainText, _ := p.list.GetItemText(idx)
		return mainText
	}
	return p.selectedKC
}

func (p *KeycodePicker) SetOnSelect(fn func(kc string)) {
	p.onSelect = fn
}

func (p *KeycodePicker) GetList() *tview.List {
	return p.list
}

func (p *KeycodePicker) GetSearchBar() *tview.InputField {
	return p.searchBar
}

func containsIgnoreCase(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	sLower := ""
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		sLower += string(c)
	}
	substrLower := ""
	for i := 0; i < len(substr); i++ {
		c := substr[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		substrLower += string(c)
	}
	return len(sLower) >= len(substrLower) && (sLower == substrLower || len(sLower) > 0 && len(substrLower) > 0 && containsSubstr(sLower, substrLower))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
