#!/bin/bash

cmd="$1"
shift

output_file=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        -o|--output)
            output_file="$2"
            shift 2
            ;;
        --output=*)
            output_file="${1#--output=}"
            shift
            ;;
        -o*)
            output_file="${1#-o}"
            shift
            ;;
        *)
            shift
            ;;
    esac
done

write_output() {
    if [[ -n "$output_file" && "$output_file" != -* ]]; then
        mkdir -p "$(dirname "$output_file")"
        cat > "$output_file"
    else
        cat
    fi
}

case "$cmd" in
    hello|list-keyboards|list-layouts|list-keymaps)
        exit 0
        ;;

    generate-version-h|generate.version-h)
        write_output <<'EOF'
#pragma once
#define QMK_VERSION "local-ci"
EOF
        exit 0
        ;;

    generate-rules-mk|generate.rules-mk)
        write_output <<'EOF'
# Generated rules.mk placeholder
EOF
        exit 0
        ;;

    generate-config-h|generate.config-h)
        write_output <<'EOF'
#pragma once
#ifndef VENDOR_ID
#define VENDOR_ID 0xFEED
#endif
#ifndef PRODUCT_ID
#define PRODUCT_ID 0x0001
#endif
#ifndef DEVICE_VER
#define DEVICE_VER 0x0001
#endif
EOF
        exit 0
        ;;

    generate-keyboard-c|generate.keyboard-c)
        write_output <<'EOF'
/* Generated keyboard C placeholder */
EOF
        exit 0
        ;;

    generate-keyboard-h|generate.keyboard-h)
        write_output <<'EOF'
#pragma once
/* Generated keyboard H placeholder */
EOF
        exit 0
        ;;

    generate-layouts|generate.layouts)
        write_output <<'EOF'
#pragma once
/* Generated layouts placeholder */
EOF
        exit 0
        ;;

    json2c)
        write_output <<'EOF'
/* Generated keymap C placeholder */
EOF
        exit 0
        ;;

    generate-dfu-header|generate.dfu-header)
        write_output <<'EOF'
/* Generated DFU header placeholder */
EOF
        exit 0
        ;;

    generate-*|generate.*)
        # Unknown generate command; create file so build can continue.
        if [[ -n "$output_file" && "$output_file" != -* ]]; then
            mkdir -p "$(dirname "$output_file")"
            touch "$output_file"
        fi
        exit 0
        ;;

    *)
        echo "qmk: unknown command: $cmd" >&2
        exit 1
        ;;
esac
