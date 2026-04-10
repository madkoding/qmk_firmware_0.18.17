# QMK ZoneKeyboards (fork basado en 0.18.17)

Este repositorio es un fork de [qmk/qmk_firmware](https://github.com/qmk/qmk_firmware) con personalizaciones de keymaps para [ZoneKeyboards](https://zonekeyboards.cl).

- Base upstream: [QMK 0.18.17](https://github.com/qmk/qmk_firmware/releases/tag/0.18.17)
- Incluye ajustes propios de build/CI para esta versión
- No corresponde a la versión más nueva de QMK oficial

## Teclados principales en este fork

- `crkbd`
- `sofle`
- `lily58`
- `lily58zk`
- `reviung41`
- `zkpad`

## Uso rápido (CLI interactivo)

Este repo incluye un CLI para compilar sin pelear con flags ni dependencias manuales.

```bash
chmod +x util/compile_cli.sh
./util/compile_cli.sh
```

El CLI permite:

- elegir un target único, la matriz CI completa o una lista custom
- compilar en modo `local` o `docker`
- preparar dependencias necesarias automáticamente
- mostrar comando de flash al final

## Requisitos

### Opción recomendada: Docker

- Docker instalado y en ejecución
- El script usa `qmkfm/qmk_cli:latest`

### Opción local (sin Docker)

- Toolchain AVR/ARM instalada
- Python + utilidades necesarias para compilar QMK

## Compilación manual (sin CLI interactivo)

Comando equivalente a CI para un target:

```bash
make <keyboard>:<keymap> SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh
```

Ejemplos:

```bash
make crkbd:zonekeyboards SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh
make lily58:zonekeyboards SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh
```

## Flash

Después de compilar, puedes flashear con:

```bash
make <keyboard>:<keymap>:flash SKIP_GIT=yes SKIP_QMK_CHECK=yes QMK_BIN=./util/qmk_wrapper.sh
```

Los binarios generados (`.hex` / `.bin`) quedan en la raíz del repo y en `.build/`.

## Notas importantes del fork

- Este repo depende de `util/qmk_wrapper.sh` para compatibilidad con la CLI usada en CI.
- Si Docker deja `.build/` con permisos `root`, puede fallar un build local posterior por permisos.
- La matriz CI actual valida: `crkbd:zonekeyboards`, `sofle:zonekeyboards`, `lily58:zonekeyboards`, `lily58zk:zonekeyboards`, `reviung41:zonekeyboards`, `zkpad:default`.
