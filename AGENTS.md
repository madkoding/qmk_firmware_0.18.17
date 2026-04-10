# Agent Notes

## Repo-specific build reality
- This fork is pinned to QMK `0.18.17`, but the Makefile still calls CLI subcommands like `list-layouts` that are missing in the container's default `qmk`.
- CI and Docker builds must use `QMK_BIN=./util/qmk_wrapper.sh` plus `SKIP_QMK_CHECK=yes` (see `.github/workflows/build-keyboards.yml`).
- `Makefile` supports overriding the CLI via `QMK_BIN ?= qmk`; do not revert this behavior.

## Canonical build commands
- CI-equivalent local build for one target:
  - `make <keyboard>:<keymap> SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh`
- Docker helper for CI matrix:
  - `./util/docker-build.sh`
  - Or pass explicit targets: `./util/docker-build.sh crkbd:zonekeyboards lily58:zonekeyboards`

## Dependency quirks (important)
- This repo does not rely on a committed `.gitmodules`; CI explicitly clones required libs into `lib/` (`lufa`, `chibios`, `chibios-contrib`, `vusb`, `printf`, `pico-sdk`).
- CI also restores `lib/lib8tion` and `lib/fnv` from upstream QMK `0.18.17` via sparse checkout.
- If a build fails on missing `lib/...` headers, check `.github/workflows/build-keyboards.yml` before changing firmware code.

## Lily58/Lily58ZK gotchas
- OLED helpers are expected under keyboard-local directories:
  - `keyboards/lily58/lib/*.c`
  - `keyboards/lily58zk/lib/*.c`
- In keymap `rules.mk`, source paths must use keyboard-local paths (for example `keyboards/lily58zk/lib/logo_reader.c`), not `./lib/...`.
- `keyboards/lily58/keymaps/zonekeyboards/rules.mk` has `LTO_ENABLE = yes` to fit AVR flash size; keep this unless you also reduce features.

## CI scope and expectations
- Current CI matrix builds only:
  - `crkbd:zonekeyboards`, `sofle:zonekeyboards`, `lily58:zonekeyboards`, `lily58zk:zonekeyboards`, `reviung41:zonekeyboards`, `zkpad:default`.
- Prefer validating changes against these targets first.

## Environment gotcha after Docker builds
- Container builds can leave root-owned `.build/` artifacts, causing local "Permission denied" on subsequent `make` runs.
- If that happens, clean ownership or remove `.build/` before retrying.
