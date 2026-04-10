package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/qmk-tui/core"
	"strings"
)

type KeyboardListView struct {
	*tview.TreeView
	keyboards []*core.Keyboard
	onSelect  func(kb *core.Keyboard, km string)
}

func NewKeyboardListView(onSelect func(kb *core.Keyboard, km string)) *KeyboardListView {
	view := &KeyboardListView{
		TreeView: tview.NewTreeView(),
		onSelect: onSelect,
	}

	view.SetTitle(" Keyboards ")
	view.SetBorder(true)
	view.SetBorderColor(tcell.NewHexColor(0xf9e2af))
	view.SetTopLevel(1)
	view.SetGraphics(true)

	keyboards, err := core.LoadKeyboards()
	if err != nil {
		panic(err)
	}
	view.keyboards = keyboards

	root := tview.NewTreeNode("QMK Keyboards (CI)")
	root.SetColor(tcell.NewHexColor(0x89b4fa))

	for _, kb := range keyboards {
		kbNode := tview.NewTreeNode(kb.Name)
		kbNode.SetReference(kb)
		kbNode.SetColor(tcell.NewHexColor(0x94e2d5))

		for _, km := range kb.Keymaps {
			displayKM := strings.TrimPrefix(km, "zonekeyboards_")
			kmNode := tview.NewTreeNode("  " + displayKM)
			kmNode.SetReference(struct {
				KB *core.Keyboard
				KM string
			}{KB: kb, KM: km})
			kmNode.SetColor(tcell.NewHexColor(0xa6e3a1))
			kbNode.AddChild(kmNode)
		}

		root.AddChild(kbNode)
	}

	view.SetRoot(root)
	root.SetExpanded(true)
	if len(root.GetChildren()) > 0 {
		view.SetCurrentNode(root.GetChildren()[0])
	}

	view.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}

		switch v := ref.(type) {
		case *core.Keyboard:
			if node.IsExpanded() {
				node.SetExpanded(false)
			} else {
				node.SetExpanded(true)
			}
		case struct {
			KB *core.Keyboard
			KM string
		}:
			if view.onSelect != nil {
				view.onSelect(v.KB, v.KM)
			}
		}
	})

	return view
}

func (v *KeyboardListView) Refresh() {
	keyboards, err := core.LoadKeyboards()
	if err != nil {
		return
	}
	v.keyboards = keyboards

	root := tview.NewTreeNode("QMK Keyboards (CI)")
	root.SetColor(tcell.NewHexColor(0x89b4fa))

	for _, kb := range keyboards {
		kbNode := tview.NewTreeNode(kb.Name)
		kbNode.SetReference(kb)
		kbNode.SetColor(tcell.NewHexColor(0x94e2d5))

		for _, km := range kb.Keymaps {
			displayKM := strings.TrimPrefix(km, "zonekeyboards_")
			kmNode := tview.NewTreeNode("  " + displayKM)
			kmNode.SetReference(struct {
				KB *core.Keyboard
				KM string
			}{KB: kb, KM: km})
			kmNode.SetColor(tcell.NewHexColor(0xa6e3a1))
			kbNode.AddChild(kmNode)
		}

		root.AddChild(kbNode)
	}

	v.SetRoot(root)
	root.SetExpanded(true)
	if len(root.GetChildren()) > 0 {
		v.SetCurrentNode(root.GetChildren()[0])
	}
}

func (v *KeyboardListView) SelectCurrentNode() {
	node := v.GetCurrentNode()
	if node == nil {
		return
	}

	ref := node.GetReference()
	if ref == nil {
		return
	}

	switch ref := ref.(type) {
	case *core.Keyboard:
		if node.IsExpanded() {
			node.SetExpanded(false)
		} else {
			node.SetExpanded(true)
		}
	case struct {
		KB *core.Keyboard
		KM string
	}:
		if v.onSelect != nil {
			v.onSelect(ref.KB, ref.KM)
		}
	}
}

func (v *KeyboardListView) MoveUp() {
	v.TreeView.Move(-1)
}

func (v *KeyboardListView) MoveDown() {
	v.TreeView.Move(1)
}
