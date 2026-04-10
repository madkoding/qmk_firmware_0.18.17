package main

import (
	"os"

	"github.com/qmk-tui/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		os.Exit(1)
	}
}
