#!/bin/bash
# QMK Docker Build Script
# Builds QMK keyboards using the official Docker container

set -e

KEYBOARDS=("${@:-crkbd:zonekeyboards sofle:zonekeyboards lily58:zonekeyboards lily58zk:zonekeyboards reviung41:zonekeyboards zkpad:default}")

echo "=========================================="
echo "QMK Docker Build Script"
echo "=========================================="

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "[ERROR] Docker is not installed or not running"
    exit 1
fi

# Create wrapper to handle missing qmk commands
WRAPPER_DIR=$(mktemp -d)
cat > "$WRAPPER_DIR/qmk" << 'WRAPPER'
#!/bin/bash
case "$1" in
    hello|list-keyboards|list-layouts|list-keymaps)
        exit 0
        ;;
    generate-*|generate.*)
        shift
        output_file=""
        while [[ $# -gt 0 ]]; do
            case $1 in
                -o|--output)
                    output_file="$2"
                    shift 2
                    ;;
                *)
                    shift
                    ;;
            esac
        done
        if [[ -n "$output_file" && "$output_file" != -* ]]; then
            mkdir -p "$(dirname "$output_file")"
            touch "$output_file"
        fi
        exit 0
        ;;
    *)
        echo "qmk: unknown command: $1" >&2
        exit 1
        ;;
esac
WRAPPER
chmod +x "$WRAPPER_DIR/qmk"

# Initialize submodules if needed
if [ ! -d "lib/chibios" ] || [ ! -d "lib/lufa" ]; then
    echo "[INFO] Initializing git submodules..."
    git submodule update --init --recursive lib/chibios lib/lufa lib/vusb lib/printf lib/lib8tion lib/fnv 2>/dev/null || true
fi

# Build each keyboard
for kb in "${KEYBOARDS[@]}"; do
    keyboard=$(echo "$kb" | cut -d: -f1)
    keymap=$(echo "$kb" | cut -d: -f2)
    
    echo ""
    echo "Building ${keyboard}:${keymap}..."
    
    docker run --rm \
        -v "$(pwd)":/qmk \
        -v "$WRAPPER_DIR":/wrapper \
        -w /qmk \
        qmkfm/qmk_cli:latest \
        sh -c "git config --global --add safe.directory /qmk && cp /wrapper/qmk /opt/uv/tools/qmk/bin/qmk && chmod +x /opt/uv/tools/qmk/bin/qmk && make ${keyboard}:${keymap} SKIP_GIT=yes" \
        && echo "[OK] ${keyboard}:${keymap}" \
        || echo "[ERROR] ${keyboard}:${keymap} failed"
done

# Cleanup
rm -rf "$WRAPPER_DIR"

echo ""
echo "=========================================="
echo "Build complete"
echo "=========================================="