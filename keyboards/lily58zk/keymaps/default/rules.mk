BOOTMAGIC_ENABLE = no       # Enable Bootmagic Lite
MOUSEKEY_ENABLE = no        # Mouse keys
EXTRAKEY_ENABLE = no        # Audio control and System control
CONSOLE_ENABLE = no         # Console for debug
COMMAND_ENABLE = no         # Commands for debug and configuration
NKRO_ENABLE = no
BACKLIGHT_ENABLE = no       # Enable keyboard backlight functionality
AUDIO_ENABLE = no           # Audio output
RGBLIGHT_ENABLE = no       # Enable WS2812 RGB underlight.
SWAP_HANDS_ENABLE = no      # Enable one-hand typing
OLED_ENABLE= yes     # OLED display

# If you want to change the display of OLED, you need to change here
SRC +=  keyboards/lily58zk/lib/rgb_state_reader.c \
        keyboards/lily58zk/lib/layer_state_reader.c \
        keyboards/lily58zk/lib/logo_reader.c \
        keyboards/lily58zk/lib/keylogger.c \
        # keyboards/lily58zk/lib/mode_icon_reader.c \
        # keyboards/lily58zk/lib/host_led_state_reader.c \
        # keyboards/lily58zk/lib/timelogger.c \
