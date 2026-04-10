# Local Build Requirements

This document describes the requirements for building QMK firmware locally.

## Required Tools

### 1. QMK CLI (>= 1.3.0)

The QMK command-line interface is required for building keyboards.

```bash
# Install from PyPI (requires Python 3.8+)
pip install qmk

# Verify installation
qmk --version
```

**Note:** The `qmk hello` command must be available for the build system to work.

### 2. AVR Toolchain

Required for compiling AVR-based keyboards (Proton C, Elite-C, etc.).

**Ubuntu/Debian:**
```bash
sudo apt-get install avr-gcc avr-libc avrdude binutils-avr
```

**Fedora:**
```bash
sudo dnf install avr-gcc avr-libc avrdude
```

**macOS:**
```bash
brew install avr-gcc avrdude
```

**Arch Linux:**
```bash
sudo pacman -S avr-gcc avr-libc avrdude
```

### 3. ARM Toolchain (optional)

Required for ARM-based keyboards (STM32, etc.).

**Ubuntu/Debian:**
```bash
sudo apt-get install gcc-arm-none-eabi libnewlib-arm-none-eabi
```

**macOS:**
```bash
brew install arm-none-eabi-gcc arm-none-eabi-binutils
```

## Git Submodules

Some libraries are included as git submodules. Initialize them with:

```bash
git submodule update --init --recursive
```

Or initialize only the required ones:
```bash
git submodule update --init lib/lufa lib/chibios lib/vusb lib/printf
```

## Building Keyboards

### Basic Build

```bash
make <keyboard>:<keymap>
```

Example:
```bash
make crkbd:zonekeyboards
make sofle:default
```

### Build with Target

```bash
make <keyboard>:<keymap>:<target>
```

Targets include:
- `:flash` - Build and flash to keyboard
- `:dfu` - Build for STM32 DFU bootloader
- `:avrdude` - Build and flash via AVRDUDE

### Parallel Build

Use `-j` to parallelize compilation:

```bash
make crkbd:zonekeyboards -j$(nproc)
```

## Using Docker

If you cannot install the toolchain locally, use the official Docker container:

```bash
docker run --rm -v $(pwd):/qmk qmkfm/qmk_cli make <keyboard>:<keymap>
```

Example:
```bash
docker run --rm -v $(pwd):/qmk qmkfm/qmk_cli make crkbd:zonekeyboards
```

## Troubleshooting

### "Cannot run qmk hello"

The QMK CLI is not properly installed or is an old version (< 1.3.0).

```bash
pip install qmk --upgrade
```

### "avr-gcc: command not found"

The AVR toolchain is not installed. See [AVR Toolchain](#avr-toolchain) above.

### "lib/lufa/LUFA/makefile: No such file or directory"

Git submodules are not initialized:

```bash
git submodule update --init lib/lufa
```

### Permission Denied on util/list_keyboards.sh

Make the script executable:

```bash
chmod +x util/list_keyboards.sh
```

## Quick Setup (Ubuntu/Debian)

```bash
# Install all requirements
sudo apt-get install -y avr-gcc avr-libc avrdude gcc-arm-none-eabi libnewlib-arm-none-eabi

# Install QMK CLI
pip install qmk --break-system-packages

# Initialize submodules
git submodule update --init --recursive

# Test build
make crkbd:zonekeyboards
```
