#!/bin/bash
# Build and run qmk_tui

cd "$(dirname "$0")"

echo "Building qmk_tui..."
docker run --rm -v $(pwd):/app -w /app golang:1.21 go build -o qmk_tui

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo ""
    echo "Running qmk_tui..."
    echo "Note: Use 'docker run' for interactive TUI with proper terminal support"
    echo ""
    ./qmk_tui
else
    echo "Build failed!"
    exit 1
fi
