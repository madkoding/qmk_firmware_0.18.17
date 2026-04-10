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