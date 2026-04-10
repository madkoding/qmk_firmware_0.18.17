package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/qmk-tui/core"
)

type KeyboardCanvas struct {
	*tview.TextView
	Layout       *core.Layout
	LayoutName   string
	Keymap       *core.Keymap
	CurrentLayer int
	SelectedIdx  int
	onKeySelect  func(key *core.Key)
}

func NewKeyboardCanvas(onKeySelect func(key *core.Key)) *KeyboardCanvas {
	canvas := &KeyboardCanvas{
		TextView:     tview.NewTextView(),
		CurrentLayer: 0,
		SelectedIdx:  -1,
		onKeySelect:  onKeySelect,
	}

	canvas.SetTitle(" Keyboard Layout ")
	canvas.SetBorder(true)
	canvas.SetBorderColor(tcell.NewHexColor(0xf9e2af))
	canvas.SetDynamicColors(true)
	canvas.SetScrollable(false)
	canvas.SetWordWrap(false)
	canvas.SetBackgroundColor(tcell.ColorBlack)

	return canvas
}

func (c *KeyboardCanvas) SetKeyboard(kb *core.Keyboard, keymap *core.Keymap) {
	c.Keymap = keymap
	c.SelectedIdx = -1

	if len(kb.Layouts) == 0 {
		return
	}

	targetLayout := kb.Layouts[0]
	if keymap != nil && keymap.LayoutName != "" {
		targetLayout = keymap.LayoutName
	}

	layoutInfo, err := kb.GetLayout(targetLayout)
	if (err != nil || layoutInfo == nil) && keymap != nil {
		if guessedName, guessedInfo := c.guessLayoutByKeyCount(kb, keymap); guessedInfo != nil {
			targetLayout = guessedName
			layoutInfo = guessedInfo
			err = nil
		}
	}
	if err != nil || layoutInfo == nil {
		layoutInfo, err = kb.GetLayout(kb.Layouts[0])
		targetLayout = kb.Layouts[0]
	}
	if err != nil || layoutInfo == nil {
		return
	}

	c.LayoutName = targetLayout
	c.Layout = core.ParseLayout(layoutInfo, targetLayout)
	c.Refresh()
}

func (c *KeyboardCanvas) guessLayoutByKeyCount(kb *core.Keyboard, keymap *core.Keymap) (string, *core.LayoutInfo) {
	if keymap == nil || len(keymap.Layers) == 0 || len(keymap.Layers[0].Keys) == 0 {
		return "", nil
	}
	target := len(keymap.Layers[0].Keys[0])
	if target == 0 {
		return "", nil
	}

	bestName := ""
	var bestInfo *core.LayoutInfo
	bestDelta := 1 << 30
	for _, name := range kb.Layouts {
		info, err := kb.GetLayout(name)
		if err != nil || info == nil {
			continue
		}
		delta := len(info.Layout) - target
		if delta < 0 {
			delta = -delta
		}
		if delta < bestDelta {
			bestDelta = delta
			bestName = name
			bestInfo = info
		}
	}

	if bestInfo != nil {
		return bestName, bestInfo
	}
	return "", nil
}

func (c *KeyboardCanvas) SetLayer(layer int) {
	if c.Layout == nil {
		return
	}
	c.CurrentLayer = layer
	c.Refresh()
}

func (c *KeyboardCanvas) Refresh() {
	if c.Layout == nil {
		c.SetText("\n\n[yellow]  Select a keyboard/keymap[white]\n[yellow]  from the list on the left[white]")
		return
	}

	var layer *core.Layer
	if c.Keymap != nil && c.CurrentLayer < len(c.Keymap.Layers) {
		layer = c.KeymapToLayer(c.Keymap, c.CurrentLayer)
	}

	// Build selected key from current layer for rendering highlight
	var selectedKey *core.Key
	if c.SelectedIdx >= 0 && layer != nil && len(layer.Keys) > 0 {
		for _, k := range layer.Keys[0] {
			if k.Index == c.SelectedIdx {
				selectedKey = k
				break
			}
		}
	}

	content := c.Layout.RenderASCIIWithSelection(layer, selectedKey)
	if c.Keymap != nil && len(c.Keymap.Layers) > 0 {
		content += "\n[#6c7086]Layers:[white] "
		for i, ly := range c.Keymap.Layers {
			if i == c.CurrentLayer {
				content += "[#89b4fa]>" + ly.Name + "<[white] "
			} else {
				content += "[#7f849c]" + ly.Name + "[white] "
			}
		}
	}
	content += "\n[#6c7086]Layout:[white] " + c.LayoutName
	if c.Keymap != nil && c.CurrentLayer < len(c.Keymap.Layers) {
		content += "   [#6c7086]Layer:[white] " + c.Keymap.Layers[c.CurrentLayer].Name
	}
	content += "\n[#6c7086]Legend:[white] " +
		"[#bac2de]Alnum/Default[white]  " +
		"[#fab387]Modifiers[white]  " +
		"[#cba6f7]Layer keys[white]  " +
		"[#f38ba8]Special[white]  " +
		"[#a6e3a1]Fn[white]  " +
		"[#94e2d5]Thumb[white]  " +
		"[#6c7086]TRNS/NO[white]"
	c.SetText(content)
}

func (c *KeyboardCanvas) KeymapToLayer(km *core.Keymap, layerIdx int) *core.Layer {
	if layerIdx >= len(km.Layers) {
		return nil
	}

	kmLayer := km.Layers[layerIdx]
	allLayoutKeys := append([]*core.Key{}, c.Layout.LeftKeys...)
	allLayoutKeys = append(allLayoutKeys, c.Layout.RightKeys...)
	allLayoutKeys = append(allLayoutKeys, c.Layout.ThumbKeys...)
	layer := &core.Layer{
		Name: kmLayer.Name,
		Keys: make([][]*core.Key, 1),
	}
	layer.Keys[0] = make([]*core.Key, 0, len(allLayoutKeys))

	if len(kmLayer.Keys) == 0 {
		return layer
	}

	flat := kmLayer.Keys[0]
	for _, lk := range allLayoutKeys {
		kc := "KC_TRNS"
		if lk.Index >= 0 && lk.Index < len(flat) && strings.TrimSpace(flat[lk.Index]) != "" {
			kc = strings.TrimSpace(flat[lk.Index])
		}
		layer.Keys[0] = append(layer.Keys[0], &core.Key{
			Index:   lk.Index,
			Keycode: kc,
			Row:     lk.Row,
			Col:     lk.Col,
			IsRight: lk.IsRight,
		})
	}

	return layer
}

func (c *KeyboardCanvas) HandleClick(x, y int) {
	if c.Layout == nil {
		return
	}

	allKeys := append([]*core.Key{}, c.Layout.LeftKeys...)
	allKeys = append(allKeys, c.Layout.RightKeys...)
	allKeys = append(allKeys, c.Layout.ThumbKeys...)

	for _, key := range allKeys {
		if x >= int(key.X) && x < int(key.X)+int(key.Width) &&
			y >= int(key.Y) && y < int(key.Y)+int(key.Height) {
			c.SelectedIdx = key.Index
			c.Refresh()
			if c.onKeySelect != nil {
				c.onKeySelect(key)
			}
			return
		}
	}
}

func (c *KeyboardCanvas) Navigate(dir string) {
	if c.Layout == nil {
		return
	}

	allKeys := append([]*core.Key{}, c.Layout.LeftKeys...)
	allKeys = append(allKeys, c.Layout.RightKeys...)
	allKeys = append(allKeys, c.Layout.ThumbKeys...)

	// If nothing selected, pick first key
	if c.SelectedIdx < 0 && len(allKeys) > 0 {
		c.SelectedIdx = allKeys[0].Index
		c.Refresh()
		if c.onKeySelect != nil {
			c.onKeySelect(allKeys[0])
		}
		return
	}

	// Find current key
	var currentKey *core.Key
	for _, k := range allKeys {
		if k.Index == c.SelectedIdx {
			currentKey = k
			break
		}
	}
	if currentKey == nil {
		return
	}

	var newKey *core.Key
	switch dir {
	case "up":
		newKey = c.findNearestKey(currentKey, 0, -1, allKeys)
	case "down":
		newKey = c.findNearestKey(currentKey, 0, 1, allKeys)
	case "left":
		newKey = c.findNearestKey(currentKey, -1, 0, allKeys)
	case "right":
		newKey = c.findNearestKey(currentKey, 1, 0, allKeys)
	}

	if newKey != nil && newKey.Index != c.SelectedIdx {
		c.SelectedIdx = newKey.Index
		c.Refresh()
		if c.onKeySelect != nil {
			c.onKeySelect(newKey)
		}
	}
}

func (c *KeyboardCanvas) findNearestKey(from *core.Key, dx, dy int, allKeys []*core.Key) *core.Key {
	var best *core.Key
	bestDist := 1000

	for _, key := range allKeys {
		if key.Index == from.Index {
			continue
		}

		dist := 0
		if dx != 0 {
			dist = abs(key.Col - from.Col)
		} else {
			dist = abs(key.Row - from.Row)
		}

		if dist < bestDist {
			bestDist = dist
			best = key
		}
	}

	return best
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (c *KeyboardCanvas) GetSelectedKey() *core.Key {
	if c.SelectedIdx < 0 || c.Layout == nil {
		return nil
	}
	allKeys := append([]*core.Key{}, c.Layout.LeftKeys...)
	allKeys = append(allKeys, c.Layout.RightKeys...)
	allKeys = append(allKeys, c.Layout.ThumbKeys...)
	for _, k := range allKeys {
		if k.Index == c.SelectedIdx {
			return k
		}
	}
	return nil
}

func (c *KeyboardCanvas) SetKeycode(kc string) {
	if c.SelectedIdx < 0 || c.Keymap == nil {
		return
	}

	if c.CurrentLayer < 0 || c.CurrentLayer >= len(c.Keymap.Layers) {
		return
	}

	if len(c.Keymap.Layers[c.CurrentLayer].Keys) == 0 {
		c.Keymap.Layers[c.CurrentLayer].Keys = [][]string{{}}
	}

	for len(c.Keymap.Layers[c.CurrentLayer].Keys[0]) <= c.SelectedIdx {
		c.Keymap.Layers[c.CurrentLayer].Keys[0] = append(c.Keymap.Layers[c.CurrentLayer].Keys[0], "KC_TRNS")
	}
	c.Keymap.Layers[c.CurrentLayer].Keys[0][c.SelectedIdx] = kc
	c.Refresh()
}
