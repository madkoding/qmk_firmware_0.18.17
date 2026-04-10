#include <stdbool.h>
#include QMK_KEYBOARD_H

extern keymap_config_t keymap_config;

#ifdef RGBLIGHT_ENABLE
#endif

#ifdef OLED_DRIVER_ENABLE
static uint32_t oled_timer = 0;
#endif

extern uint8_t is_master;

typedef union {
  struct {
    bool enable_upper_light :1;
    bool enable_downer_light :1;
  };
} user_config_t;
user_config_t user_config;

int upper_keys[6] = {9, 10, 17, 18, 23, 24};

enum layers {
  _QWERTY,
  _LOWER,
  _HIGHER,
  _ADJUST
};

enum custom_keycodes {
  QWERTY = SAFE_RANGE,
  LOWER,
  HIGHER,
  ADJUST,
  TOOGLE_DOWNER,
  TOOGLE_UPPER,
};

// Unicode
enum unicode_names {
    a_CODE,
    A_CODE,
    e_CODE,
    E_CODE,
    i_CODE,
    I_CODE,
    o_CODE,
    O_CODE,
    u_CODE,
    U_CODE,
};

const uint32_t PROGMEM unicode_map[] = {
    [a_CODE]  = 0x00E1,
    [A_CODE]  = 0x00C1,
    [e_CODE]  = 0x00E9,
    [E_CODE]  = 0x00C9,
    [i_CODE]  = 0x00ED,
    [I_CODE]  = 0x00CD,
    [o_CODE]  = 0x00F3,
    [O_CODE]  = 0x00D3,
    [u_CODE]  = 0x00FA,
    [U_CODE]  = 0x00DA,
};

// Tap Dance definitions
enum {
    TD_CAPLOCK,
    TD_SLASH,
    TD_LEFT_PAR_BRA,
    TD_RIGHT_PAR_BRA,
    TD_A,
    TD_E,
    TD_I,
    TD_O,
    TD_U,
};

qk_tap_dance_action_t tap_dance_actions[] = {
    [TD_CAPLOCK] = ACTION_TAP_DANCE_DOUBLE(KC_LSFT, KC_CAPS),
    [TD_SLASH] = ACTION_TAP_DANCE_DOUBLE(KC_PSLS, A(KC_NUBS)),
    [TD_LEFT_PAR_BRA] = ACTION_TAP_DANCE_DOUBLE(A(KC_LBRC), A(KC_QUOT)),
    [TD_RIGHT_PAR_BRA] = ACTION_TAP_DANCE_DOUBLE(A(KC_RBRC), A(KC_BSLS)),
    [TD_A] = ACTION_TAP_DANCE_DOUBLE(KC_A, XP(a_CODE, A_CODE)),
    [TD_E] = ACTION_TAP_DANCE_DOUBLE(KC_E, XP(e_CODE, E_CODE)),
    [TD_I] = ACTION_TAP_DANCE_DOUBLE(KC_I, XP(i_CODE, I_CODE)),
    [TD_O] = ACTION_TAP_DANCE_DOUBLE(KC_O, XP(o_CODE, O_CODE)),
    [TD_U] = ACTION_TAP_DANCE_DOUBLE(KC_U, XP(u_CODE, U_CODE)),
};

const uint16_t PROGMEM keymaps[][MATRIX_ROWS][MATRIX_COLS] = {[_QWERTY] = LAYOUT(
  //|-----------------------------------------------------|                    |-----------------------------------------------------|
     KC_ESC,  KC_Q,    KC_W,    KC_E,    KC_R,    KC_T,                         KC_Y,    KC_U,    KC_I,    KC_O,    KC_P,    KC_BSPC,
  //|--------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     KC_TAB,  KC_A,    KC_S,    KC_D,    KC_F,    KC_G,                         KC_H,    KC_J,    KC_K,    KC_L,    KC_SCLN, KC_LCTL,
  //---------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     TD(TD_CAPLOCK),KC_Z,KC_X,  KC_C,    KC_V,    KC_B,                         KC_N,    KC_M,    KC_COMM, KC_DOT,  KC_GRV,  KC_QUOT,
  //---------+--------+--------+--------+--------+--------+--------|  |--------+--------+--------+--------+--------+--------+--------|
                                         KC_LGUI, LOWER, KC_SPC,       KC_ENT,  HIGHER,  KC_LALT
                                      //|--------------------------|  |--------------------------|
  ),

  [_LOWER] = LAYOUT(
  //|-----------------------------------------------------|                    |-----------------------------------------------------|
     KC_GRV,  KC_1,    KC_2,    KC_3,    KC_4,    KC_5,                         KC_6,    KC_7,    KC_8,    KC_9,    KC_0,    KC_BSPC,
  //|--------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     KC_TAB,  KC_MINS, S(KC_2), S(KC_3), A(KC_3), S(KC_5),                      KC_PAST, KC_PPLS, KC_SLSH, TD(TD_SLASH), KC_UP, KC_DEL,
  //|--------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     KC_LSFT, KC_EQL,  S(KC_1), A(KC_2), S(KC_8), S(KC_9),                      TD(TD_LEFT_PAR_BRA),TD(TD_RIGHT_PAR_BRA), S(KC_0), KC_LEFT, KC_DOWN, KC_RGHT,
  //|--------+--------+--------+--------+--------+--------+--------|  |--------+--------+--------+--------+--------+--------+--------|
                                        KC_LGUI, LOWER,   KC_SPC,     KC_ENT,  HIGHER,  KC_RALT
                                      //|--------------------------|  |--------------------------|
  ),

  [_HIGHER] = LAYOUT(
  //|-----------------------------------------------------|                    |-----------------------------------------------------|
     KC_ESC,  A(KC_1), KC_LCBR, S(KC_QUOT), A(KC_SCLN), S(KC_NUBS),             KC_BSLS, A(KC_E), XXXXXXX, XXXXXXX, XXXXXXX, KC_BSPC,
  //|--------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     KC_F1,   KC_F2,   KC_F3,   KC_F4,   KC_F5,   KC_F6,                        XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX,
  //|--------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     KC_F7,   KC_F8,   KC_F9,   KC_F10,  KC_F11,  KC_F12,                       XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX,
  //|--------+--------+--------+--------+--------+--------+--------|  |--------+--------+--------+--------+--------+--------+--------|
                                         KC_LCTL, LOWER,   KC_SPC,     KC_ENT,  HIGHER,  KC_RALT
                                      //|--------------------------|  |--------------------------|
  ),

  [_ADJUST] = LAYOUT(
  //|-----------------------------------------------------|                    |-----------------------------------------------------|
     KC_ESC,  RGB_VAD, RGB_VAI, RGB_MOD,RGB_TOG, TOOGLE_DOWNER,                 TOOGLE_UPPER, RGB_SPI,  RGB_SPD, XXXXXXX, XXXXXXX, KC_BSPC,
  //|--------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     XXXXXXX, RGB_HUD, RGB_HUI, KC__MUTE, KC__VOLDOWN, KC__VOLUP,               KC_MPRV, KC_MPLY,  KC_MNXT, XXXXXXX, XXXXXXX, XXXXXXX,
  //|--------+--------+--------+--------+--------+--------|                    |--------+--------+--------+--------+--------+--------|
     XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX,                      XXXXXXX, XXXXXXX,  XXXXXXX, XXXXXXX, XXXXXXX, RESET,
  //|--------+--------+--------+--------+--------+--------+--------|  |--------+--------+--------+--------+--------+--------+--------|
                                         KC_LCTL, LOWER,   KC_SPC,     KC_ENT,  HIGHER,  KC_RALT
                                      //|--------------------------|  |--------------------------|
  )};

int RGB_current_mode;

// Setting ADJUST layer RGB back to default
void update_tri_layer_RGB(uint8_t layer1, uint8_t layer2, uint8_t layer3) {
  if (IS_LAYER_ON(layer1) && IS_LAYER_ON(layer2)) {
    layer_on(layer3);
  } else {
    layer_off(layer3);
  }
}

void matrix_init_user(void) {
    #ifdef RGBLIGHT_ENABLE
      RGB_current_mode = rgblight_config.mode;
    #endif
    user_config.enable_upper_light = true;
    user_config.enable_downer_light = true;
}

void rgb_matrix_indicators_user(void) {
  #ifdef RGB_MATRIX_ENABLE
  switch (biton32(layer_state)) {
    case _HIGHER:
        if (user_config.enable_upper_light) {
            for (int i = 0; i < 6; i++) {
            rgb_matrix_set_color(upper_keys[i], 80, 80, 80);
            }
        }
        if (user_config.enable_downer_light) {
            rgb_matrix_set_color(13, 255,83,0);
            rgb_matrix_set_color(14, 255,83,0);
            rgb_matrix_set_color(6, 255,83,0);
        }
      break;
    case _LOWER:
        if (user_config.enable_upper_light) {
            for (int i = 0; i < 6; i++) {
            rgb_matrix_set_color(upper_keys[i], 80, 80, 80);
            }
        }
        if (user_config.enable_downer_light) {
            rgb_matrix_set_color(13, 153,204,0);
            rgb_matrix_set_color(14, 153,204,0);
            rgb_matrix_set_color(6, 153,204,0);
        }
      break;
    case _ADJUST:
        if (user_config.enable_upper_light) {
            for (int i = 0; i < 6; i++) {
            rgb_matrix_set_color(upper_keys[i], 80, 80, 80);
            }
        }
        if (user_config.enable_downer_light) {
            rgb_matrix_set_color(13, 204,0,0);
            rgb_matrix_set_color(14, 204,0,0);
            rgb_matrix_set_color(6, 204,0,0);
        }
      break;
    default:
        if (host_keyboard_leds() & (1<<USB_LED_CAPS_LOCK)) {
            rgb_matrix_set_color(13, 179,255,255);
            rgb_matrix_set_color(14, 179,255,255);
            rgb_matrix_set_color(6, 179,255,255);
            rgb_matrix_set_color(26, 179,255,255);
        }
        if (user_config.enable_upper_light) {
            for (int i = 0; i < 6; i++) {
            rgb_matrix_set_color(upper_keys[i], 80, 80, 80);
            }
        }
      break;
  }
  #endif
}

#ifdef OLED_DRIVER_ENABLE
oled_rotation_t oled_init_user(oled_rotation_t rotation) { return OLED_ROTATION_270; }

void render_space(void) {
    oled_write_P(PSTR("     "), false);
}

void render_mod_status_gui_alt(uint8_t modifiers) {
    static const char PROGMEM gui_off_1[] = {0x85, 0x86, 0};
    static const char PROGMEM gui_off_2[] = {0xa5, 0xa6, 0};
    static const char PROGMEM gui_on_1[] = {0x8d, 0x8e, 0};
    static const char PROGMEM gui_on_2[] = {0xad, 0xae, 0};

    static const char PROGMEM alt_off_1[] = {0x87, 0x88, 0};
    static const char PROGMEM alt_off_2[] = {0xa7, 0xa8, 0};
    static const char PROGMEM alt_on_1[] = {0x8f, 0x90, 0};
    static const char PROGMEM alt_on_2[] = {0xaf, 0xb0, 0};

    // fillers between the modifier icons bleed into the icon frames
    static const char PROGMEM off_off_1[] = {0xc5, 0};
    static const char PROGMEM off_off_2[] = {0xc6, 0};
    static const char PROGMEM on_off_1[] = {0xc7, 0};
    static const char PROGMEM on_off_2[] = {0xc8, 0};
    static const char PROGMEM off_on_1[] = {0xc9, 0};
    static const char PROGMEM off_on_2[] = {0xca, 0};
    static const char PROGMEM on_on_1[] = {0xcb, 0};
    static const char PROGMEM on_on_2[] = {0xcc, 0};

    if(modifiers & MOD_MASK_GUI) {
        oled_write_P(gui_on_1, false);
    } else {
        oled_write_P(gui_off_1, false);
    }

    if ((modifiers & MOD_MASK_GUI) && (modifiers & MOD_MASK_ALT)) {
        oled_write_P(on_on_1, false);
    } else if(modifiers & MOD_MASK_GUI) {
        oled_write_P(on_off_1, false);
    } else if(modifiers & MOD_MASK_ALT) {
        oled_write_P(off_on_1, false);
    } else {
        oled_write_P(off_off_1, false);
    }

    if(modifiers & MOD_MASK_ALT) {
        oled_write_P(alt_on_1, false);
    } else {
        oled_write_P(alt_off_1, false);
    }

    if(modifiers & MOD_MASK_GUI) {
        oled_write_P(gui_on_2, false);
    } else {
        oled_write_P(gui_off_2, false);
    }

    if (modifiers & MOD_MASK_GUI & MOD_MASK_ALT) {
        oled_write_P(on_on_2, false);
    } else if(modifiers & MOD_MASK_GUI) {
        oled_write_P(on_off_2, false);
    } else if(modifiers & MOD_MASK_ALT) {
        oled_write_P(off_on_2, false);
    } else {
        oled_write_P(off_off_2, false);
    }

    if(modifiers & MOD_MASK_ALT) {
        oled_write_P(alt_on_2, false);
    } else {
        oled_write_P(alt_off_2, false);
    }
}

void render_mod_status_ctrl_shift(uint8_t modifiers) {
    static const char PROGMEM ctrl_off_1[] = {0x89, 0x8a, 0};
    static const char PROGMEM ctrl_off_2[] = {0xa9, 0xaa, 0};
    static const char PROGMEM ctrl_on_1[] = {0x91, 0x92, 0};
    static const char PROGMEM ctrl_on_2[] = {0xb1, 0xb2, 0};

    static const char PROGMEM shift_off_1[] = {0x8b, 0x8c, 0};
    static const char PROGMEM shift_off_2[] = {0xab, 0xac, 0};
    static const char PROGMEM shift_on_1[] = {0xcd, 0xce, 0};
    static const char PROGMEM shift_on_2[] = {0xcf, 0xd0, 0};

    // fillers between the modifier icons bleed into the icon frames
    static const char PROGMEM off_off_1[] = {0xc5, 0};
    static const char PROGMEM off_off_2[] = {0xc6, 0};
    static const char PROGMEM on_off_1[] = {0xc7, 0};
    static const char PROGMEM on_off_2[] = {0xc8, 0};
    static const char PROGMEM off_on_1[] = {0xc9, 0};
    static const char PROGMEM off_on_2[] = {0xca, 0};
    static const char PROGMEM on_on_1[] = {0xcb, 0};
    static const char PROGMEM on_on_2[] = {0xcc, 0};

    if(modifiers & MOD_MASK_CTRL) {
        oled_write_P(ctrl_on_1, false);
    } else {
        oled_write_P(ctrl_off_1, false);
    }

    if ((modifiers & MOD_MASK_CTRL) && (modifiers & MOD_MASK_SHIFT)) {
        oled_write_P(on_on_1, false);
    } else if(modifiers & MOD_MASK_CTRL) {
        oled_write_P(on_off_1, false);
    } else if(modifiers & MOD_MASK_SHIFT) {
        oled_write_P(off_on_1, false);
    } else {
        oled_write_P(off_off_1, false);
    }

    if(modifiers & MOD_MASK_SHIFT) {
        oled_write_P(shift_on_1, false);
    } else {
        oled_write_P(shift_off_1, false);
    }

    if(modifiers & MOD_MASK_CTRL) {
        oled_write_P(ctrl_on_2, false);
    } else {
        oled_write_P(ctrl_off_2, false);
    }

    if (modifiers & MOD_MASK_CTRL & MOD_MASK_SHIFT) {
        oled_write_P(on_on_2, false);
    } else if(modifiers & MOD_MASK_CTRL) {
        oled_write_P(on_off_2, false);
    } else if(modifiers & MOD_MASK_SHIFT) {
        oled_write_P(off_on_2, false);
    } else {
        oled_write_P(off_off_2, false);
    }

    if(modifiers & MOD_MASK_SHIFT) {
        oled_write_P(shift_on_2, false);
    } else {
        oled_write_P(shift_off_2, false);
    }
}

void render_logo(void) {
    static const char PROGMEM corne_logo[] = {
        0x80, 0x81, 0x82, 0x83, 0x84,
        0xa0, 0xa1, 0xa2, 0xa3, 0xa4,
        0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0};
    oled_write_P(corne_logo, false);
    oled_write_P(PSTR("corne"), false);
}

void render_layer_state(void) {
    static const char PROGMEM default_layer[] = {
        0x20, 0x94, 0x95, 0x96, 0x20,
        0x20, 0xb4, 0xb5, 0xb6, 0x20,
        0x20, 0xd4, 0xd5, 0xd6, 0x20, 0};
    static const char PROGMEM raise_layer[] = {
        0x20, 0x97, 0x98, 0x99, 0x20,
        0x20, 0xb7, 0xb8, 0xb9, 0x20,
        0x20, 0xd7, 0xd8, 0xd9, 0x20, 0};
    static const char PROGMEM lower_layer[] = {
        0x20, 0x9a, 0x9b, 0x9c, 0x20,
        0x20, 0xba, 0xbb, 0xbc, 0x20,
        0x20, 0xda, 0xdb, 0xdc, 0x20, 0};
    static const char PROGMEM adjust_layer[] = {
        0x20, 0x9d, 0x9e, 0x9f, 0x20,
        0x20, 0xbd, 0xbe, 0xbf, 0x20,
        0x20, 0xdd, 0xde, 0xdf, 0x20, 0};
    if(layer_state_is(_ADJUST)) {
        oled_write_P(adjust_layer, false);
    } else if(layer_state_is(_LOWER)) {
        oled_write_P(lower_layer, false);
    } else if(layer_state_is(_HIGHER)) {
        oled_write_P(raise_layer, false);
    } else {
        oled_write_P(default_layer, false);
    }
}

void render_status_main(void) {
    render_logo();
    render_space();
    render_layer_state();
    render_space();
    render_mod_status_gui_alt(get_mods()|get_oneshot_mods());
    render_mod_status_ctrl_shift(get_mods()|get_oneshot_mods());
}

void render_status_secondary(void) {
    render_logo();
    render_space();
    render_layer_state();
    render_space();
    render_mod_status_gui_alt(get_mods()|get_oneshot_mods());
    render_mod_status_ctrl_shift(get_mods()|get_oneshot_mods());
}

void oled_task_user(void) {
    if (timer_elapsed32(oled_timer) > 1500000) {
        oled_off();
        return;
    }
#ifndef SPLIT_KEYBOARD
    else { oled_on(); }
#endif

    if (is_master) {
        render_status_main();  // Renders the current keyboard state (layer, lock, caps, scroll, etc)
    } else {
        render_status_secondary();
    }
}
#endif // OLED_DRIVER_ENABLE

 bool process_record_user(uint16_t keycode, keyrecord_t *record) {
  if (record->event.pressed) {
    #ifdef OLED_DRIVER_ENABLE
            oled_timer = timer_read32();
    #endif
    // set_timelog();
  }

  switch (keycode) {
    case LOWER:
      if (record->event.pressed) {
        layer_on(_LOWER);
        update_tri_layer_RGB(_LOWER, _HIGHER, _ADJUST);
      } else {
        layer_off(_LOWER);
        update_tri_layer_RGB(_LOWER, _HIGHER, _ADJUST);
      }
      return false;
    case HIGHER:
      if (record->event.pressed) {
        layer_on(_HIGHER);
        update_tri_layer_RGB(_LOWER, _HIGHER, _ADJUST);
      } else {
        layer_off(_HIGHER);
        update_tri_layer_RGB(_LOWER, _HIGHER, _ADJUST);
      }
      return false;
    case ADJUST:
        if (record->event.pressed) {
          layer_on(_ADJUST);
        } else {
          layer_off(_ADJUST);
        }
        return false;
    case TOOGLE_DOWNER:
        if (record->event.pressed) {
            user_config.enable_downer_light = !user_config.enable_downer_light;
        }
        return false;
    case TOOGLE_UPPER:
        if (record->event.pressed) {
            user_config.enable_upper_light = !user_config.enable_upper_light;
        }
        return false;
  }
  return true;
}

#ifdef RGB_MATRIX_ENABLE

    void suspend_power_down_user(void) {
        rgb_matrix_set_suspend_state(true);
    }

    void suspend_wakeup_init_user(void) {
        rgb_matrix_set_suspend_state(false);
    }

#endif
