/* Copyright 2020 gtips
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
#include QMK_KEYBOARD_H


enum keycodes {
    QWERTY = SAFE_RANGE,
    LOWER,
    RAISE,
    ADJUST
};

enum layer_names {
    _QWERTY,
    _LOWER,
    _RAISE,
    _ADJUST
};

// Tap Dance definitions
enum {
    TD_CAPLOCK,
    TD_SLASH,
    TD_LEFT_PAR_BRA,
    TD_RIGHT_PAR_BRA,
};

qk_tap_dance_action_t tap_dance_actions[] = {
    [TD_CAPLOCK] = ACTION_TAP_DANCE_DOUBLE(KC_LSFT, KC_CAPS),
    [TD_SLASH] = ACTION_TAP_DANCE_DOUBLE(KC_PSLS, A(KC_NUBS)),
    [TD_LEFT_PAR_BRA] = ACTION_TAP_DANCE_DOUBLE(A(KC_LBRC), A(KC_QUOT)),
    [TD_RIGHT_PAR_BRA] = ACTION_TAP_DANCE_DOUBLE(A(KC_RBRC), A(KC_BSLS))
};

const uint16_t PROGMEM keymaps[][MATRIX_ROWS][MATRIX_COLS] = {

/* Qwerty
 * ,-----------------------------------------.     ,-----------------------------------------.
 * | Tab  |   Q  |   W  |   E  |   R  |   T  |     |   Y  |   U  |   I  |   O  |   P  |BackSp|
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * |LCTRL |   A  |   S  |   D  |   F  |   G  |     |   H  |   J  |   K  |   L  |   Ñ  | ´¨{  |
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * |LShift|   Z  |   X  |   C  |   V  |   B  |     |   N  |   M  |  ,;  |  .:  |  -_  |Enter |
 * `-----------------------------------------/     \-----------------------------------------'
 *                          |Command| LOWER| / Space \ |RAISE |  Alt  |
 *                          |       |      |/         \|      |       |
 *                          `------------------------''--------------'
 */

  [_QWERTY] = LAYOUT_reviung41(
    KC_TAB,          KC_Q,     KC_W,     KC_E,     KC_R,      KC_T,               KC_Y,     KC_U,     KC_I,     KC_O,     KC_P,     KC_BSPC,
    KC_LCTL,         KC_A,     KC_S,     KC_D,     KC_F,      KC_G,               KC_H,     KC_J,     KC_K,     KC_L,     KC_SCLN,  KC_QUOT,
    TD(TD_CAPLOCK),  KC_Z,     KC_X,     KC_C,     KC_V,      KC_B,               KC_N,     KC_M,     KC_COMM,  KC_DOT,   KC_SLSH,  RSFT_T(KC_ENT),
                                                   KC_LGUI,   LOWER,    KC_SPC,   RAISE,    KC_LALT
  ),

/* Lower
 * ,-----------------------------------------.     ,-----------------------------------------.
 * | ESC  |  ºª\ |  ¿   |   ?  |   ^  |  *   |     |  7/÷ |  8(“ |  9)” |   *  |  /   |BackSP|
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * |LShift|  <>  |  [{  |  ]}  |   '  |  "   |     |  4$¢ |  5%∞ |  6&¬ |   +  |  -   |  C   |
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * |LCTRL |  (   |  )   |  /\  | `^[  | +*]  |     |  1!| |  2"@ | 3·#  |  0=≠ |  .   |  =   |
 * `-----------------------------------------/     \-----------------------------------------'
 *                          |Command| LOWER| / Enter \ |RAISE |  Alt  |
 *                          |       |      |/         \|      |       |
 *                          `------------------------''--------------'
 */
  [_LOWER] = LAYOUT_reviung41(
    KC_ESCAPE,  KC_NUBS,  KC_UNDS,  KC_PLUS,  KC_LCBR,   KC_RCBR,                        KC_7,     KC_8,     KC_9,     KC_ASTR,  KC_PSLS, _______,
    _______,    KC_GRV,TD(TD_LEFT_PAR_BRA),TD(TD_RIGHT_PAR_BRA),KC_MINS,S(KC_2),         KC_4,     KC_5,     KC_6,     KC_PLUS,  KC_PMNS, S(KC_C),
    _______,    S(KC_8), S(KC_9), TD(TD_SLASH), KC_LBRC, KC_RBRC,                        KC_1,     KC_2,     KC_3,     KC_0,     KC_DOT,  KC_PEQL,
                                            _______,   _______,  KC_ENT,   _______,  _______
  ),


/* Raise
 * ,-----------------------------------------.     ,-----------------------------------------.
 * | Esc  |  |   |  ^   |   ̈   |  ~   |  a   |     |      |      |  Up  |Pag UP|Pag DW|BackSP|
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * |      |      |      |      |      |      |     | Home | Left | Down |Right | END  |  DEL |
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * | F1   | F2   | F3   | F4   | F5   | F6   |     |  F7  |  F8  |  F9  |  F10 |  F11 |  F12 |
 * `-----------------------------------------/     \-----------------------------------------'
 *                          |Command| LOWER| / Space \ |RAISE |  Alt  |
 *                          |       |      |/         \|      |       |
 *                          `------------------------''--------------'
 */
  [_RAISE] = LAYOUT_reviung41(
    KC_ESC,  A(KC_1), KC_LCBR, S(KC_QUOT), A(KC_SCLN), S(KC_NUBS),             XXXXXXX, XXXXXXX, KC_UP, KC_PGUP, KC_PGDN, _______,
    XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX, XXXXXXX,                    KC_HOME, KC_LEFT, KC_DOWN, KC_RGHT, KC_END,  KC_DEL,
    KC_F1,   KC_F2,   KC_F3,   KC_F4,   KC_F5,   KC_F6,                      KC_F7,   KC_F8,   KC_F9,   KC_F10, KC_F11, KC_F12,
                                            _______,   _______,  _______,  _______,  _______
  ),


/* Adjust
 * ,-----------------------------------------.     ,-----------------------------------------.
 * |      |      |      |      |      |RGBTog|     |      |      |      | Mute | VOL- | VOL+ |
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * |      | Luz+ | Sat+ | Mat+ | Spd+ | Mod  |     |      |      |      |      |      |      |
 * |------+------+------+------+------+------|     |------+------+------+------+------+------|
 * |      | Luz- | Sat- | Mat- | Spd- |      |     | Reset|      |      |      |      |      |
 * `-----------------------------------------/     \-----------------------------------------'
 *                          |Command| LOWER| /      \ |RAISE |  Alt  |
 *                          |      |      |/         \|      |       |
 *                          `----------------------------------------'
 */
  [_ADJUST] = LAYOUT_reviung41(
    XXXXXXX,   XXXXXXX, XXXXXXX,  XXXXXXX,  XXXXXXX,   RGB_TOG,            XXXXXXX,  XXXXXXX,  XXXXXXX,  KC__MUTE, KC__VOLDOWN, KC__VOLUP,
    XXXXXXX,   RGB_HUI, RGB_SAI,  RGB_VAI,  RGB_SPI,   RGB_MOD,            XXXXXXX,  XXXXXXX,  XXXXXXX,  XXXXXXX,  XXXXXXX,     XXXXXXX,
    XXXXXXX,   RGB_HUD, RGB_SAD,  RGB_VAD,  RGB_SPD,    XXXXXXX,           RESET,    XXXXXXX,  XXXXXXX,  XXXXXXX,  XXXXXXX,     XXXXXXX,
                                            _______,   _______,  XXXXXXX,  _______,  _______
  )
};

void update_tri_layer_RGB(uint8_t layer1, uint8_t layer2, uint8_t layer3) {
    if (IS_LAYER_ON(layer1) && IS_LAYER_ON(layer2)) {
        rgblight_sethsv_at(HSV_PURPLE, 10);
        layer_on(layer3);
    } else {
        // rgblight_sethsv_at(HSV_GREEN, 10);
        layer_off(layer3);
    }
}

bool process_record_user(uint16_t keycode, keyrecord_t *record){
    switch (keycode) {
        case LOWER:
            if (record->event.pressed) {
                layer_on(_LOWER);
                rgblight_sethsv_at(HSV_BLUE, 10);
                update_tri_layer_RGB(_LOWER, _RAISE, _ADJUST);
            } else {
                layer_off(_LOWER);
                rgblight_sethsv_at(HSV_GREEN, 10);
                update_tri_layer_RGB(_LOWER, _RAISE, _ADJUST);
            }
            return false;
        case RAISE:
            if (record->event.pressed) {
                layer_on(_RAISE);
                rgblight_sethsv_at(HSV_RED, 10);
                update_tri_layer_RGB(_LOWER, _RAISE, _ADJUST);
            } else {
                layer_off(_RAISE);
                rgblight_sethsv_at(HSV_GREEN, 10);
                update_tri_layer_RGB(_LOWER, _RAISE, _ADJUST);
            }
            return false;
        case ADJUST:
            if (record->event.pressed) {
                layer_on(_ADJUST);
                rgblight_sethsv_at(HSV_PURPLE, 10);
                update_tri_layer_RGB(_LOWER, _RAISE, _ADJUST);
            } else {
                layer_off(_ADJUST);
                rgblight_sethsv_at(HSV_GREEN, 10);
                update_tri_layer_RGB(_LOWER, _RAISE, _ADJUST);
            }
            return false;
            break;
    }
    return true;
}

void keyboard_post_init_user(void) {
    rgblight_sethsv_at(HSV_GREEN, 10);
    rgblight_set_effect_range( 0, 10);
}
