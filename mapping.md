# Remap Action Reference — DrunkDeer CLI

## Type 0 — Standard Keys

### Letters

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| A | A | N | N |
| B | B | O | O |
| C | C | P | P |
| D | D | Q | Q |
| E | E | R | R |
| F | F | S | S |
| G | G | T | T |
| H | H | U | U |
| I | I | V | V |
| J | J | W | W |
| K | K | X | X |
| L | L | Y | Y |
| M | M | Z | Z |

### Numerals

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| 0 | 0 | 5 | 5 |
| 1 | 1 | 6 | 6 |
| 2 | 2 | 7 | 7 |
| 3 | 3 | 8 | 8 |
| 4 | 4 | 9 | 9 |

### Symbols & System

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| BACK | BACK | ESC | ESC |
| CAPS | CAPS | SPACE | SPACE |
| ENTER | RETURN | TAB | TAB |
| MINUS | MINUS | PLUS | PLUS |
| COMMA | COMMA | PERIOD | PERIOD |
| COLON | COLON | QOTATN | QOTATN |
| SLASH | SLASH | SLASH_K29 | SLASH_K29 |
| BRKTS_L | BRKTS_L | BRKTS_R | BRKTS_R |
| EUR_K45 | --- | K45 | --- |
| TILDE | --- | | |

### F-keys & Navigation (none on G60)

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| F1 | --- | F7 | --- |
| F2 | --- | F8 | --- |
| F3 | --- | F9 | --- |
| F4 | --- | F10 | --- |
| F5 | --- | F11 | --- |
| F6 | --- | F12 | --- |
| PRINT | --- | SCRLK | --- |
| PAUSE | --- | INS | --- |
| HOME | --- | PGUP | --- |
| DEL | --- | END | --- |
| PGDN | --- | ARR_L | --- |
| ARR_R | --- | ARR_UP | --- |
| ARR_DW | --- | APP | --- |
| LOOP | --- | NOLOOP | --- |
| KANA | --- | K56 | --- |

### Numpad (none on G60)

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| NUMS | --- | NUM0 | --- |
| NUM1 | --- | NUM2 | --- |
| NUM3 | --- | NUM4 | --- |
| NUM5 | --- | NUM6 | --- |
| NUM7 | --- | NUM8 | --- |
| NUM9 | --- | KP_DEL | --- |
| KP_ENTER | --- | KP_PLUS | --- |
| KP_MINUS | --- | KP_MULT | --- |
| KP_DIV | --- | | |

## Type 1 — Modifiers

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| CTRL_L | CTRL_L | CTRL_R | CTRL_R |
| SHIFT_L / SHF_L | SHF_L | SHIFT_R / SHF_R | SHF_R |
| ALT_L | ALT_L | ALT_R | ALT_R |
| WIN_L | WIN_L | SCRLK | --- |
| FN1 | FN1 | FN2 | FN2 |

## Type 2 — Features (none on G60 by default)

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| WIN_LOCK | --- | RDT | --- |
| LW | --- | RT | --- |
| RT_MATCH | --- | RT_EXTREME | --- |
| RT_STANDARD | --- | RT_COMPETE | --- |
| LIGHT_SW | --- | LIGHT_CYCLE | --- |
| LIGHT_NEXT | --- | LIGHT_PREV | --- |
| LIGHT_BR_INC | --- | LIGHT_BR_DEC | --- |
| LIGHT_SP_INC | --- | LIGHT_SP_DEC | --- |
| LIGHT_COLOR | --- | LIGHT_PAUSE | --- |
| LIGHT_CTRL | --- | COLOR_CTRL | --- |
| SPEED_CTRL | --- | BR_CTRL | --- |

## Type 3 — Multimedia & Mouse (none on G60 by default)

| Action | Default G60 key | Action | Default G60 key |
|--------|-----------------|--------|-----------------|
| VOL_UP | --- | VOL_DN | --- |
| MUTE | --- | PLAY | --- |
| STOP | --- | PREV | --- |
| NEXT | --- | MS_L | --- |
| MS_R | --- | MS_M | --- |
| MS_SCR_U | --- | MS_SCR_D | --- |
| MS_SCR_L | --- | MS_SCR_R | --- |

---

# G60 Physical Key Layout (61 keys)

### Row 0 (top)
ESC  1  2  3  4  5  6  7  8  9  0  MINUS  PLUS  BACK

### Row 1
TAB  Q  W  E  R  T  Y  U  I  O  P  BRKTS_L  BRKTS_R  SLASH_K29

### Row 2
CAPS  A  S  D  F  G  H  J  K  L  COLON  QOTATN  RETURN

### Row 3
SHF_L  Z  X  C  V  B  N  M  COMMA  PERIOD  SLASH  SHF_R

### Row 4
CTRL_L  WIN_L  ALT_L  SPACE  ALT_R  FN1  FN2  CTRL_R

# Full KEYBOARD_LAYOUT (generic 104-key, indices 0–125)

```
  0: ESC       21: TILDE     42: TAB       63: CAPS      84: SHF_L     105: CTRL_L
  1: (empty)   22: 1         43: Q         64: A         85: (empty)   106: WIN_L
  2: F1        23: 2         44: W         65: S         86: Z         107: ALT_L
  3: F2        24: 3         45: E         66: D         87: X         108: (empty)
  4: F3        25: 4         46: R         67: F         88: C         109: (empty)
  5: F4        26: 5         47: T         68: G         89: V         110: (empty)
  6: F5        27: 6         48: Y         69: H         90: B         111: SPACE
  7: F6        28: 7         49: U         70: I         91: N         112: (empty)
  8: F7        29: 8         50: I         71: J         92: M         113: (empty)
  9: F8        30: 9         51: O         72: K         93: COMMA     114: (empty)
 10: F9        31: 0         52: P         73: L         94: PERIOD    115: ALT_R
 11: F10       32: MINUS     53: BRKTS_L   74: COLON     95: SLASH     116: FN1
 12: F11       33: PLUS      54: BRKTS_R   75: (empty)   96: (empty)   117: APP
 13: F12       34: BACK      55: SLASH_K29 76: RETURN    97: SHF_R     118: ARR_L
 14: NUM7       35: (empty)   56: (empty)   77: (empty)   98: (empty)   119: ARR_DW
 15: NUM8       36: HOME      57: PGUP      78: (empty)   99: END       120: ARR_R
 16: NUM9       37: (empty)   58: (empty)   79: (empty)  100: (empty)   121: CTRL_R
 17: (empty)   38: (empty)   59: (empty)   80: (empty)  101: (empty)   122: (empty)
 18: (empty)   39: (empty)   60: (empty)   81: (empty)  102: (empty)   123: (empty)
 19: (empty)   40: (empty)   61: (empty)   82: (empty)  103: (empty)   124: (empty)
 20: (empty)   41: (empty)   62: (empty)   83: (empty)  104: (empty)   125: (empty)
```

# Aliases

| Name | Alias of | Same firmware value |
|------|----------|-------------------|
| SHF_L | SHIFT_L | `{0xFC, 2, 1}` |
| SHF_R | SHIFT_R | `{0xFC, 32, 1}` |
| K45 | EUR_K45 | `{0xFC, 0x64, 0}` |
| RETURN | ENTER | `{0xFC, 0x28, 0}` |
