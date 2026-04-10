# QMK Firmware FIXME Tracking

This document tracks all documented FIXME comments in the codebase.
See https://github.com/qmk/qmk_firmware/issues for existing issues.

## Summary

| Category | Files | Count |
|----------|-------|-------|
| quantum/ | ~15 | 105+ |
| tmk_core/ | ~5 | 25+ |
| users/ | ~10 | 15+ |
| **Total** | **~30** | **140+** |

---

## Critical Files with High FIXME Density

### quantum/eeconfig.c (24 FIXME)
Functions needing documentation:
- eeconfig_init (line 23)
- eeconfig_read_* functions (lines 38, 91, 99, 107, 118, 132, 146, 153, 161, 168, 176, 183, 192, 199, 207, 214, 222, 229, 237, 244, 252, 259)

### quantum/action_util.c (26 FIXME)
- Functions: default_layer_state_, set_single_paste_default_layer, update_reference, suspend_state, waking_state, debug_config_set, debug_config_get, default_layer_set, layer_state_set, autoshift_enabled_, autoshift_disable_, tap_toggle_, wait, waiting_buffer_, waiting_buffer_send_, waiting_buffer_ended_, waiting_buffer_flush_, housekeeping_task_, debug_code_, debug_action_, debug_event_, debug_state_, and more.

### quantum/keyboard.c (15 FIXME)
- Functions: matrix_init_* callbacks, matrix_scan_* callbacks, proc_keycode, default_layer_* functions

### quantum/backlight/backlight.c (14 FIXME)
- Functions: backlight_* (init, set, breathing*, effects, etc.)

### quantum/action_tapping.c (8 FIXME)
- Functions: debug_tapping, pre_modifier, after_modifier, reprocess_tapping, waiting_buffer_, qk_action_tapping_

---

## Documentation Improvements Needed

### quantum/main.c
- `main()` (line 48) - FIXME: Needs doc

### quantum/action_layer.c
- `action_for_key()` (line 298) - FIXME: Needs better summary

### quantum/action.c
- Multiple action_* functions need documentation (lines 76, 172, 228, 252, 311)

### quantum/bootmagic/magic.c
- bootmagic_* functions (line 34)

---

## TODO vs FIXME

This repo uses FIXME for missing documentation and implementation gaps.
Consider converting common FIXME patterns to issues:
- https://github.com/qmk/qmk_firmware/labels/documentation
- https://github.com/qmk/qmk_firmware/labels/improvement

---

## How to Help

1. Pick a file from this list
2. Look at similar documented functions in the codebase
3. Follow the existing docstring style
4. Submit a PR removing the FIXME and adding proper documentation

## Documentation Style Example

```c
/**
 * \brief Description of what the function does
 *
 * More detailed description if needed,
 * explaining edge cases or behavior.
 *
 * \param param_name Description of parameter
 * \return Description of return value
 */
void function_name(param_type param_name) {
    // implementation
}
```

---

Last updated: 2026-04-09
Generated from: grep -r "FIXME" quantum/ tmk_core/ users/ --include="*.c"
