package core

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Key struct {
	Index    int
	Label    string
	X, Y     float64
	Width    float64
	Height   float64
	Keycode  string
	Layer    int
	Row      int
	Col      int
	IsRight  bool
	IsThumb  bool
}

type Layer struct {
	Name string
	Keys [][]*Key
}

type Layout struct {
	Name         string
	LayoutMacro  string
	LeftKeys     []*Key
	RightKeys    []*Key
	ThumbKeys    []*Key
	Rows         int
	Cols         int
	TotalKeys    int
}

func ParseLayout(info *LayoutInfo, layoutName string) *Layout {
	layout := &Layout{
		Name:        layoutName,
		LayoutMacro: layoutName,
	}

	var leftKeys []*Key
	var rightKeys []*Key
	var thumbKeys []*Key
	var maxY float64

	for i := range info.Layout {
		keyInfo := &info.Layout[i]
		key := &Key{
			Index:   i,
			Label:   keyInfo.Label,
			X:       keyInfo.X,
			Y:       keyInfo.Y,
			Width:   keyInfo.Width,
			Height:  keyInfo.Height,
			Keycode: "",
			Layer:   0,
			Row:     int(keyInfo.Y),
			Col:     int(keyInfo.X),
		}

		if key.Width == 0 {
			key.Width = 1
		}
		if key.Height == 0 {
			key.Height = 1
		}

		if key.Y > maxY {
			maxY = key.Y
		}

		isRight := key.X >= estimateSplitX(info.Layout)
		key.IsRight = isRight

		if key.Y >= 3 {
			key.IsThumb = true
			thumbKeys = append(thumbKeys, key)
		} else if isRight {
			rightKeys = append(rightKeys, key)
		} else {
			leftKeys = append(leftKeys, key)
		}

	}

	layout.LeftKeys = leftKeys
	layout.RightKeys = rightKeys
	layout.ThumbKeys = thumbKeys
	layout.Rows = int(maxY) + 1
	layout.TotalKeys = len(info.Layout)

	return layout
}

func estimateSplitX(keys []KeyInfo) float64 {
	if len(keys) == 0 {
		return 6
	}
	xs := make([]float64, 0, len(keys))
	for _, k := range keys {
		xs = append(xs, k.X)
	}
	sort.Float64s(xs)

	maxGap := -1.0
	split := xs[len(xs)/2]
	for i := 1; i < len(xs); i++ {
		gap := xs[i] - xs[i-1]
		if gap > maxGap {
			maxGap = gap
			split = (xs[i] + xs[i-1]) / 2
		}
	}
	if maxGap < 1.2 {
		return split
	}
	return split
}

func (l *Layout) RenderASCII(layer *Layer) string {
	return l.RenderASCIIWithSelection(layer, nil)
}

func (l *Layout) RenderASCIIWithSelection(layer *Layer, selected *Key) string {
	var b strings.Builder

	b.WriteString("[#cdd6f4:#000000]")
	b.WriteString("\n")
	b.WriteString("[#f9e2af]╔══════════════════════════════════════════════════════════════════╗[white]\n")
	b.WriteString("[#f9e2af]║[#cdd6f4]                 [#89b4fa]QMK KEYBOARD LAYOUT[#cdd6f4]                              [#f9e2af]║[white]\n")
	b.WriteString("[#f9e2af]╚══════════════════════════════════════════════════════════════════╝[white]\n\n")

	keycodeMap := make(map[string]string)
	if layer != nil {
		for _, row := range layer.Keys {
			for _, k := range row {
				if k != nil {
					posKey := fmt.Sprintf("idx:%d", k.Index)
					keycodeMap[posKey] = k.Keycode
				}
			}
		}
	}
	b.WriteString("[#6c7086]PHYSICAL MAP[white]\n")
	allKeys := append([]*Key{}, l.LeftKeys...)
	allKeys = append(allKeys, l.RightKeys...)
	allKeys = append(allKeys, l.ThumbKeys...)
	b.WriteString(l.renderGeometric(allKeys, keycodeMap, selected))

	return b.String()
}

func (l *Layout) renderGeometric(keys []*Key, keycodeMap map[string]string, selected *Key) string {
	var b strings.Builder
	if len(keys) == 0 {
		return "  (none)\n"
	}

	rowBuckets := map[int][]*Key{}
	rows := make([]int, 0)
	for _, k := range keys {
		ry := int(math.Round(k.Y * 2))
		if _, ok := rowBuckets[ry]; !ok {
			rows = append(rows, ry)
		}
		rowBuckets[ry] = append(rowBuckets[ry], k)
	}
	sort.Ints(rows)

	for _, row := range rows {
		rowKeys := rowBuckets[row]
		sort.Slice(rowKeys, func(i, j int) bool { return rowKeys[i].X < rowKeys[j].X })

		top := strings.Builder{}
		mid := strings.Builder{}
		bot := strings.Builder{}
		top.WriteString("  ")
		mid.WriteString("  ")
		bot.WriteString("  ")
		cursor := 2
		for _, k := range rowKeys {
			target := 2 + int(math.Round(k.X*8.0))
			if target > cursor {
				pad := strings.Repeat(" ", target-cursor)
				top.WriteString(pad)
				mid.WriteString(pad)
				bot.WriteString(pad)
				cursor = target
			}

			t, m, bo, w := l.renderKeySquare(k, keycodeMap, selected)
			top.WriteString(t)
			mid.WriteString(m)
			bot.WriteString(bo)
			cursor += w
		}
		b.WriteString(top.String())
		b.WriteString("\n")
		b.WriteString(mid.String())
		b.WriteString("\n")
		b.WriteString(bot.String())
		b.WriteString("\n")
	}

	return b.String()
}

func (l *Layout) renderKeySquare(k *Key, keycodeMap map[string]string, selected *Key) (string, string, string, int) {
	kc := keycodeMap[fmt.Sprintf("idx:%d", k.Index)]
	if kc == "" {
		kc = "____"
	}

	label := formatKeyLabel(kc, k.Label)
	width := int(math.Round(k.Width * 6.0))
	if width < 6 {
		width = 6
	}
	inner := width - 2
	if len(label) > inner {
		label = label[:inner]
	}

	isSelected := selected != nil && k.Index == selected.Index

	color := getKeyColor(kc, k.IsThumb)

	if isSelected {
		top := fmt.Sprintf("[#89b4fa]┏%s┓[white]", strings.Repeat("━", inner))
		mid := fmt.Sprintf("[#1e1e2e:#89b4fa]┃%s┃[white]", padCenter(label, inner))
		bot := fmt.Sprintf("[#89b4fa]┗%s┛[white]", strings.Repeat("━", inner))
		return top, mid, bot, width
	}
	top := fmt.Sprintf("%s┌%s┐[white]", color, strings.Repeat("─", inner))
	mid := fmt.Sprintf("%s│%s│[white]", color, padCenter(label, inner))
	bot := fmt.Sprintf("%s└%s┘[white]", color, strings.Repeat("─", inner))
	return top, mid, bot, width
}

func formatKeyLabel(kc, fallback string) string {
	k := strings.TrimSpace(kc)
	if k == "" {
		k = fallback
	}

	if strings.HasPrefix(k, "KC_") {
		raw := strings.TrimPrefix(k, "KC_")
		switch raw {
		case "BSPC":
			return "BSP"
		case "ENTER", "ENT":
			return "ENT"
		case "ESC", "ESCAPE":
			return "ESC"
		case "SPACE", "SPC":
			return "SPC"
		case "TAB":
			return "TAB"
		case "LEFT":
			return "LFT"
		case "RIGHT":
			return "RGT"
		case "DOWN":
			return "DWN"
		case "UP":
			return "UP"
		case "TRNS":
			return "TRNS"
		case "NO":
			return "NO"
		}
		if len(raw) <= 4 {
			return raw
		}
		return raw[:4]
	}

	if strings.HasPrefix(k, "MO(") {
		return "MO"
	}
	if strings.HasPrefix(k, "TG(") {
		return "TG"
	}
	if strings.HasPrefix(k, "LT(") {
		return "LT"
	}
	if strings.HasPrefix(k, "MT(") {
		return "MT"
	}
	if strings.HasPrefix(k, "DF(") {
		return "DF"
	}

	if len(k) <= 4 {
		return k
	}
	if fallback != "" && len(fallback) <= 4 {
		return fallback
	}
	return k[:4]
}

func getKeyColor(kc string, isThumb bool) string {
	if isThumb {
		return "[#94e2d5]"
	}

	switch {
	case strings.HasPrefix(kc, "KC_L") || strings.HasPrefix(kc, "KC_R"):
		return "[#fab387]"
	case strings.HasPrefix(kc, "MO(") || strings.HasPrefix(kc, "TG(") || strings.HasPrefix(kc, "LT(") || strings.HasPrefix(kc, "DF("):
		return "[#cba6f7]"
	case strings.HasPrefix(kc, "MT("):
		return "[#fab387]"
	case kc == "KC_ESCAPE", kc == "KC_ENTER", kc == "KC_SPACE", kc == "KC_TAB", kc == "KC_BSPC":
		return "[#f38ba8]"
	case strings.HasPrefix(kc, "KC_F") && len(kc) <= 5:
		return "[#a6e3a1]"
	case kc == "KC_NO" || kc == "KC_TRNS":
		return "[#6c7086]"
	default:
		return "[#bac2de]"
	}
}

func padCenter(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := width - len(s)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func (l *Layout) GetKeyAt(x, y float64) *Key {
	allKeys := append([]*Key{}, l.LeftKeys...)
	allKeys = append(allKeys, l.RightKeys...)
	allKeys = append(allKeys, l.ThumbKeys...)

	for _, k := range allKeys {
		if x >= k.X && x < k.X+k.Width && y >= k.Y && y < k.Y+k.Height {
			return k
		}
	}
	return nil
}

func (l *Layout) FindKey(row, col int, isRight bool) *Key {
	allKeys := append([]*Key{}, l.LeftKeys...)
	allKeys = append(allKeys, l.RightKeys...)
	allKeys = append(allKeys, l.ThumbKeys...)

	for _, k := range allKeys {
		if k.Row == row && k.Col == col && k.IsRight == isRight {
			return k
		}
	}
	return nil
}
