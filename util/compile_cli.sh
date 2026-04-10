#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

TARGETS=(
  "crkbd:zonekeyboards"
  "sofle:zonekeyboards"
  "lily58:zonekeyboards"
  "lily58zk:zonekeyboards"
  "reviung41:zonekeyboards"
  "zkpad:default"
)

SELECTED_TARGETS=()

ensure_repo_root() {
  if [[ ! -f "$ROOT_DIR/Makefile" ]]; then
    echo "[ERROR] No se encontro Makefile en $ROOT_DIR"
    exit 1
  fi
}

ensure_permissions() {
  chmod +x "$ROOT_DIR/util/list_keyboards.sh" "$ROOT_DIR/util/qmk_wrapper.sh"
}

ensure_lib_deps() {
  mkdir -p "$ROOT_DIR/lib"

  [[ -d "$ROOT_DIR/lib/lufa/.git" ]] || git clone --depth 1 https://github.com/qmk/lufa.git "$ROOT_DIR/lib/lufa"
  [[ -d "$ROOT_DIR/lib/chibios/.git" ]] || git clone --depth 1 https://github.com/qmk/ChibiOS.git "$ROOT_DIR/lib/chibios"
  [[ -d "$ROOT_DIR/lib/chibios-contrib/.git" ]] || git clone --depth 1 https://github.com/qmk/ChibiOS-Contrib.git "$ROOT_DIR/lib/chibios-contrib"
  [[ -d "$ROOT_DIR/lib/vusb/.git" ]] || git clone --depth 1 https://github.com/qmk/v-usb.git "$ROOT_DIR/lib/vusb"
  [[ -d "$ROOT_DIR/lib/printf/.git" ]] || git clone --depth 1 https://github.com/qmk/printf.git "$ROOT_DIR/lib/printf"
  [[ -d "$ROOT_DIR/lib/pico-sdk/.git" ]] || git clone --depth 1 https://github.com/qmk/pico-sdk.git "$ROOT_DIR/lib/pico-sdk"

  if [[ ! -f "$ROOT_DIR/lib/lib8tion/lib8tion.h" || ! -f "$ROOT_DIR/lib/fnv/fnv.h" ]]; then
    rm -rf /tmp/qmk_upstream_cli
    git clone --depth 1 --filter=blob:none --sparse --branch 0.18.17 https://github.com/qmk/qmk_firmware.git /tmp/qmk_upstream_cli
    git -C /tmp/qmk_upstream_cli sparse-checkout set lib/lib8tion lib/fnv
    cp -R /tmp/qmk_upstream_cli/lib/lib8tion "$ROOT_DIR/lib/lib8tion"
    cp -R /tmp/qmk_upstream_cli/lib/fnv "$ROOT_DIR/lib/fnv"
  fi
}

pick_single_target() {
  echo ""
  echo "Selecciona un target para compilar:"
  local i=1
  for t in "${TARGETS[@]}"; do
    printf "  %d) %s\n" "$i" "$t"
    i=$((i + 1))
  done
  printf "  %d) custom\n" "$i"

  echo ""
  read -r -p "Opcion: " opt

  if [[ "$opt" =~ ^[0-9]+$ ]] && (( opt >= 1 && opt <= ${#TARGETS[@]} )); then
    SELECTED_TARGETS=("${TARGETS[$((opt - 1))]}")
    return
  fi

  if [[ "$opt" == "$i" ]]; then
    read -r -p "Ingresa target (ej. crkbd:zonekeyboards): " custom
    SELECTED_TARGETS=("$custom")
    return
  fi

  echo "[ERROR] Opcion invalida"
  exit 1
}

pick_targets() {
  echo ""
  echo "Que deseas compilar?"
  echo "  1) un target"
  echo "  2) matriz CI completa (${#TARGETS[@]} targets)"
  echo "  3) lista custom (separada por espacios)"
  echo ""
  read -r -p "Opcion: " scope

  case "$scope" in
    1)
      pick_single_target
      ;;
    2)
      SELECTED_TARGETS=("${TARGETS[@]}")
      ;;
    3)
      read -r -p "Ingresa targets (ej. crkbd:zonekeyboards lily58:zonekeyboards): " line
      read -r -a SELECTED_TARGETS <<< "$line"
      if [[ ${#SELECTED_TARGETS[@]} -eq 0 ]]; then
        echo "[ERROR] No ingresaste targets"
        exit 1
      fi
      ;;
    *)
      echo "[ERROR] Opcion invalida"
      exit 1
      ;;
  esac
}

pick_mode() {
  echo ""
  echo "Modo de compilacion:"
  echo "  1) local (toolchain instalado)"
  echo "  2) docker (recomendado)"
  echo ""
  read -r -p "Opcion: " mode

  case "$mode" in
    1) MODE="local" ;;
    2) MODE="docker" ;;
    *)
      echo "[ERROR] Opcion invalida"
      exit 1
      ;;
  esac
}

run_local_build() {
  local target="$1"
  (cd "$ROOT_DIR" && make "$target" SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh)
}

run_docker_build() {
  local target="$1"
  if ! command -v docker >/dev/null 2>&1; then
    echo "[ERROR] Docker no esta disponible"
    exit 1
  fi

  docker run --rm \
    -v "$ROOT_DIR":/qmk \
    -w /qmk \
    qmkfm/qmk_cli:latest \
    sh -c "chmod +x util/list_keyboards.sh util/qmk_wrapper.sh && make $target SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh" < /dev/null
}

show_flash_hint() {
  local target="$1"
  local kb="${target%%:*}"
  local km="${target##*:}"
  echo ""
  echo "Compilacion completada para $target"
  echo ""
  echo "Para flashear manualmente:"
  echo "  make ${kb}:${km}:flash SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh"
}

ask_flash_now() {
  echo ""
  echo "Deseas flashear ahora?"
  echo "  1) no"
  echo "  2) si (solo local)"
  echo ""
  read -r -p "Opcion: " do_flash

  case "$do_flash" in
    1)
      return
      ;;
    2)
      if [[ "$MODE" != "local" ]]; then
        echo "[WARN] El flasheo directo solo esta soportado en modo local."
        echo "       Ejecuta localmente: make <keyboard>:<keymap>:flash SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh"
        return
      fi

      for target in "${SELECTED_TARGETS[@]}"; do
        local kb="${target%%:*}"
        local km="${target##*:}"
        echo ""
        echo "Flasheando $target ..."
        (cd "$ROOT_DIR" && make "${kb}:${km}:flash" SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh)
      done
      ;;
    *)
      echo "[ERROR] Opcion invalida"
      exit 1
      ;;
  esac
}

main() {
  ensure_repo_root
  ensure_permissions
  pick_targets
  pick_mode

  echo ""
  echo "Preparando dependencias..."
  ensure_lib_deps

  echo ""
  for target in "${SELECTED_TARGETS[@]}"; do
    echo ""
    echo "Compilando $target en modo $MODE..."
    if [[ "$MODE" == "local" ]]; then
      run_local_build "$target"
    else
      run_docker_build "$target"
    fi
    show_flash_hint "$target"
  done

  echo ""
  echo "Ultimos .hex generados:"
  (cd "$ROOT_DIR" && ls -1t *.hex 2>/dev/null | head -n 10) || true

  ask_flash_now
}

main
