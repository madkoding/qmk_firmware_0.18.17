#!/bin/bash
# Run QMK TUI - requires interactive terminal

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
QMK_TUI_DIR="$SCRIPT_DIR/qmk_tui"
QMK_BINARY="$QMK_TUI_DIR/qmk_tui"

if [ ! -f "$QMK_BINARY" ]; then
    echo "Building qmk_tui..."
    cd "$QMK_TUI_DIR"
    docker run --rm -v "$(pwd)":/app -w /app golang:1.21 go build -o qmk_tui
    cd "$SCRIPT_DIR"
fi

echo ""
echo "Starting QMK TUI..."
echo "Controls: j/k/h/l navigate | Enter edit | 1-4 layers | b=build | f=flash | s=save | Ctrl+C exit"
echo ""

cd "$SCRIPT_DIR"

exec docker run -it --rm \
    -v "$SCRIPT_DIR":/qmk \
    -e TERM=xterm-256color \
    -e QMK_PATH=/qmk \
    -w /qmk \
    --privileged \
    qmkfm/qmk_cli:latest \
    /qmk/qmk_tui/qmk_tui
