# QMK TUI - Keyboard Layout Manager

Interfaz TUI para compilar, flashear y editar keymaps de teclados QMK.

## Características

- **Visualización de layouts**: Muestra la distribución física de teclas del teclado
- **Edición de keymaps**: Click o navegación con flechas para seleccionar y cambiar teclas
- **Compilación y flash**: Build integrado con `make` y flasheo automático
- **Copy layout**: Copia la distribución de otro keymap del mismo teclado
- **Navegación por teclado**: vim-like (h,j,k,l) + atajos para layers

## Teclados Soportados (CI)

- crkbd (Corne)
- lily58 / lily58zk
- sofle
- reviung41
- zkpad

## Controles

| Tecla | Acción |
|-------|--------|
| `j/k/h/l` | Navegar entre teclas (vim-style) |
| `Enter` | Editar tecla seleccionada |
| `1-4` | Cambiar layer |
| `b` | Build |
| `f` | Flash |
| `s` | Guardar |
| `Ctrl+C` | Salir |

## Requisitos

- Go 1.21+ (para compilar)
- Docker (alternativa para ejecutar)
- QMK firmware en el path actual o `QMK_PATH`

## Build

```bash
# Con Go instalado localmente
cd qmk_tui
go build -o qmk_tui

# Con Docker
cd ..
docker run --rm -v $(pwd):/app -w /app golang:1.21 go build -o qmk_tui/qmk_tui qmk_tui
```

## Ejecución

```bash
# Local (necesita terminal interactiva)
./qmk_tui/qmk_tui

# Con Docker (recomendado)
docker run -it --rm \
  -v $(pwd)/..:/qmk \
  -e TERM=xterm-256color \
  -e QMK_PATH=/qmk \
  qmkfm/qmk_cli \
  /qmk_tui/qmk_tui

# O construir y ejecutar con script
cd qmk_tui && ./build_and_run.sh
```

## Notas

- El path de QMK se detecta automáticamente buscando en `.` o usando `QMK_PATH`
- Los cambios se guardan directamente en `keyboards/<kb>/keymaps/<km>/keymap.c`
- El parser de keymaps funciona con formatos típicos de QMK (puede no funcionar con macros complejos)
