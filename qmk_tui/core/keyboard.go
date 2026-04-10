package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Keyboard struct {
	Name       string
	Path       string
	Layouts    []string
	Bootloader string
	MCU        string
	Keymaps    []string
}

type LayoutInfo struct {
	Layout []KeyInfo `json:"layout"`
}

type KeyInfo struct {
	Label  string  `json:"label"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"w"`
	Height float64 `json:"h"`
}

type InfoJSON struct {
	KeyboardName string                  `json:"keyboard_name"`
	Layouts       map[string]LayoutInfo `json:"layouts"`
}

var CI_KEYBOARDS = []string{
	"crkbd",
	"lily58",
	"lily58zk",
	"sofle",
	"reviung41",
	"zkpad",
}

func GetQMKPath() string {
	if qmkPath := os.Getenv("QMK_PATH"); qmkPath != "" {
		return qmkPath
	}
	return "."
}

func LoadKeyboards() ([]*Keyboard, error) {
	var keyboards []*Keyboard
	qmkPath := GetQMKPath()

	for _, name := range CI_KEYBOARDS {
		kb, err := loadKeyboard(qmkPath, name)
		if err != nil {
			continue
		}
		keyboards = append(keyboards, kb)
	}

	sort.Slice(keyboards, func(i, j int) bool {
		return keyboards[i].Name < keyboards[j].Name
	})

	return keyboards, nil
}

func loadKeyboard(qmkPath, name string) (*Keyboard, error) {
	kbPath := filepath.Join(qmkPath, "keyboards", name)

	infoJSONPath := filepath.Join(kbPath, "info.json")
	if _, err := os.Stat(infoJSONPath); os.IsNotExist(err) {
		revisions, _ := os.ReadDir(kbPath)
		for _, rev := range revisions {
			if rev.IsDir() && rev.Name() != "keymaps" && rev.Name() != "lib" {
				infoJSONPath = filepath.Join(kbPath, rev.Name(), "info.json")
				if _, err := os.Stat(infoJSONPath); err == nil {
					break
				}
			}
		}
	}

	data, err := os.ReadFile(infoJSONPath)
	if err != nil {
		return nil, err
	}

	var info InfoJSON
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	rulesPath := filepath.Join(kbPath, "rules.mk")
	rulesData, _ := os.ReadFile(rulesPath)
	rulesContent := string(rulesData)

	kb := &Keyboard{
		Name:       name,
		Path:       kbPath,
		Bootloader: extractValue(rulesContent, "BOOTLOADER"),
		MCU:        extractValue(rulesContent, "MCU"),
	}

	for layoutName := range info.Layouts {
		kb.Layouts = append(kb.Layouts, layoutName)
	}
	sort.Strings(kb.Layouts)

	keymapsPath := filepath.Join(kbPath, "keymaps")
	if entries, err := os.ReadDir(keymapsPath); err == nil {
		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				kb.Keymaps = append(kb.Keymaps, entry.Name())
			}
		}
	}
	sort.Strings(kb.Keymaps)

	return kb, nil
}

func extractValue(content, key string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, key+" ") || strings.HasPrefix(line, key+"\t") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func (k *Keyboard) GetLayout(layoutName string) (*LayoutInfo, error) {
	infoJSONPath := filepath.Join(k.Path, "info.json")
	if _, err := os.Stat(infoJSONPath); os.IsNotExist(err) {
		revisions, _ := os.ReadDir(k.Path)
		for _, rev := range revisions {
			if rev.IsDir() && !strings.HasPrefix(rev.Name(), ".") {
				infoJSONPath = filepath.Join(k.Path, rev.Name(), "info.json")
				if _, err := os.Stat(infoJSONPath); err == nil {
					break
				}
			}
		}
	}

	data, err := os.ReadFile(infoJSONPath)
	if err != nil {
		return nil, err
	}

	var info InfoJSON
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	if layout, ok := info.Layouts[layoutName]; ok {
		return &layout, nil
	}

	return nil, nil
}
