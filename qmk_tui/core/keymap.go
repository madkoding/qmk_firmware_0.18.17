package core

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Keymap struct {
	Keyboard   *Keyboard
	KeymapName string
	LayoutName string
	Layers     []*KeymapLayer
	FilePath   string
}

type KeymapLayer struct {
	Name  string
	Index int
	Macro string
	Keys  [][]string
}

func LoadKeymap(kb *Keyboard, keymapName string) (*Keymap, error) {
	keymapPath := filepath.Join(kb.Path, "keymaps", keymapName, "keymap.c")

	data, err := os.ReadFile(keymapPath)
	if err != nil {
		return nil, err
	}

	km := &Keymap{
		Keyboard:   kb,
		KeymapName: keymapName,
		LayoutName: "LAYOUT",
		FilePath:   keymapPath,
	}

	km.Layers = parseKeymapContent(string(data), kb)
	if len(km.Layers) > 0 && km.Layers[0].Macro != "" {
		km.LayoutName = km.Layers[0].Macro
	}

	return km, nil
}

func parseKeymapContent(content string, kb *Keyboard) []*KeymapLayer {
	var layers []*KeymapLayer

	layerPattern := regexp.MustCompile(`\[(\w+)\]\s*=\s*(LAYOUT[\w]*)\s*\(`)
	matches := layerPattern.FindAllStringSubmatchIndex(content, -1)

	if len(matches) == 0 {
		return layers
	}

	for i, match := range matches {
		start := match[0]
		end := len(content)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}

		layerContent := content[start:end]

		layerName := "Layer"
		nameMatch := regexp.MustCompile(`\[(\w+)\]`).FindStringSubmatch(layerContent)
		if len(nameMatch) > 1 {
			layerName = nameMatch[1]
		}

		macro := "LAYOUT"
		macroMatch := regexp.MustCompile(`\[(\w+)\]\s*=\s*(LAYOUT[\w]*)\s*\(`).FindStringSubmatch(layerContent)
		if len(macroMatch) > 2 {
			macro = macroMatch[2]
		}

		keys := extractLayerKeys(layerContent, kb)

		layers = append(layers, &KeymapLayer{
			Name:  layerName,
			Index: i,
			Macro: macro,
			Keys:  keys,
		})
	}

	return layers
}

func extractLayerKeys(content string, kb *Keyboard) [][]string {
	start := strings.Index(content, "LAYOUT")
	if start == -1 {
		return nil
	}
	open := strings.Index(content[start:], "(")
	if open == -1 {
		return nil
	}
	open = start + open

	depth := 0
	close := -1
	for i := open; i < len(content); i++ {
		switch content[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				close = i
				break
			}
		}
	}
	if close == -1 || close <= open+1 {
		return nil
	}

	args := content[open+1 : close]
	tokens := splitKeyTokens(args)
	if len(tokens) == 0 {
		return nil
	}
	return [][]string{tokens}
}

func splitKeyTokens(args string) []string {
	var tokens []string
	var current strings.Builder
	depth := 0

	flush := func() {
		t := strings.TrimSpace(current.String())
		if t != "" {
			tokens = append(tokens, t)
		}
		current.Reset()
	}

	for i := 0; i < len(args); i++ {
		ch := args[i]
		switch ch {
		case '(':
			depth++
			current.WriteByte(ch)
		case ')':
			if depth > 0 {
				depth--
			}
			current.WriteByte(ch)
		case ',':
			if depth == 0 {
				flush()
			} else {
				current.WriteByte(ch)
			}
		case '\n', '\r', '\t':
			current.WriteByte(' ')
		default:
			current.WriteByte(ch)
		}
	}
	flush()
	return tokens
}

func (km *Keymap) Save() error {
	content := generateKeymapC(km)
	return os.WriteFile(km.FilePath, []byte(content), 0644)
}

func generateKeymapC(km *Keymap) string {
	var b strings.Builder

	b.WriteString("#include QMK_KEYBOARD_H\n\n")

	if len(km.Layers) > 0 {
		b.WriteString("enum layers {\n")
		for i, layer := range km.Layers {
			b.WriteString("\t" + layer.Name)
			if i < len(km.Layers)-1 {
				b.WriteString(",\n")
			} else {
				b.WriteString("\n")
			}
		}
		b.WriteString("};\n\n")
	}

	b.WriteString("const uint16_t PROGMEM keymaps[][MATRIX_ROWS][MATRIX_COLS] = {\n\n")
	for _, layer := range km.Layers {
		macro := km.LayoutName
		if layer.Macro != "" {
			macro = layer.Macro
		}
		b.WriteString("\t[" + layer.Name + "] = " + macro + "(\n")

		for rowIdx, row := range layer.Keys {
			b.WriteString("\t\t")
			for colIdx, key := range row {
				b.WriteString(key)
				if colIdx < len(row)-1 {
					b.WriteString(", ")
				}
			}
			b.WriteString(",")
			if rowIdx < len(layer.Keys)-1 {
				b.WriteString("\n")
			} else {
				b.WriteString(" // row\n")
			}
		}

		b.WriteString("\t),\n\n")
	}
	b.WriteString("};\n\n")

	b.WriteString("bool process_record_user(uint16_t keycode, keyrecord_t *record) {\n")
	b.WriteString("\tswitch (keycode) {\n")
	b.WriteString("\tdefault:\n")
	b.WriteString("\t\treturn true;\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn false;\n")
	b.WriteString("}\n")

	return b.String()
}

func CopyKeymap(src *Keymap, destKB *Keyboard, destName string) (*Keymap, error) {
	newKM := &Keymap{
		Keyboard:   destKB,
		KeymapName: destName,
		Layers:     make([]*KeymapLayer, len(src.Layers)),
		FilePath:   filepath.Join(destKB.Path, "keymaps", destName, "keymap.c"),
	}

	for i, srcLayer := range src.Layers {
		newLayer := &KeymapLayer{
			Name:  srcLayer.Name,
			Index: srcLayer.Index,
			Keys:  make([][]string, len(srcLayer.Keys)),
		}
		for j, row := range srcLayer.Keys {
			newLayer.Keys[j] = make([]string, len(row))
			copy(newLayer.Keys[j], row)
		}
		newKM.Layers[i] = newLayer
	}

	destDir := filepath.Join(destKB.Path, "keymaps", destName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}

	if err := newKM.Save(); err != nil {
		return nil, err
	}

	return newKM, nil
}

func (km *Keymap) SetKey(layerIdx, row, col int, keycode string) {
	if layerIdx < 0 || layerIdx >= len(km.Layers) {
		return
	}
	layer := km.Layers[layerIdx]
	if row < 0 || row >= len(layer.Keys) {
		return
	}
	if col < 0 || col >= len(layer.Keys[row]) {
		return
	}
	layer.Keys[row][col] = keycode
}
