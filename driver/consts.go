package driver

const (
	KEYBOARD_REPORT_ID   = 0x04
	KEYS_PER_ROW         = 59   // According to their dumb layout
	DEFAULT_ACTUATION    = 0x14 // 0x14 is 2.0mm, there should be always 59 keys!!! (unless 3rd row)
	CUSTOM_COLOR_PADDING = 0x80
	COLORS_PER_PACKET    = 13
	BYTES_PER_KEY        = 4
	G60_LED_GRID_SIZE    = 61  // 61 physical keys on G60
)

const (
	PACKET_IDENTITY    = 0xA0
	PACKET_LEDMODESEL  = 0xAE
	PACKET_MODIFYKEY   = 0xB6
	PACKET_TURBORT     = 0xB5
	PACKET_KEYTRACKING = 0xB7
)

const (
	KEYBOARD_A75PRO = "A75" // We don't give a damn about pro, they work pretty much the same
	KEYBOARD_A75    = "A75"
	KEYBOARD_G65    = "G65"
	KEYBOARD_G60    = "G60"
	KEYBOARD_G75    = "G75"
)

// G60-specific LED mode sequence values
// Derived from the webdriver's adjusted array indices (Math.max(0, arrayIndex - 2))
const (
	SEQUENCE_OFF         = 0x00
	SEQUENCE_WAVE        = 0x02 // 光谱波浪 (spectrum wave)
	SEQUENCE_SURF_RIGHT  = 0x03 // 往右冲浪 (right surfing)
	SEQUENCE_BREATH      = 0x04 // 呼吸 (breathing)
	SEQUENCE_SURF_CENTER = 0x05 // 往中冲浪 (center surfing)
	SEQUENCE_CYCLE       = 0x06 // 光谱循环 (spectrum cycle)
	SEQUENCE_RIPPLE      = 0x07 // 按键涟漪 (key ripple)
	SEQUENCE_ALWAYS      = 0x08 // 常亮 (always on)
	SEQUENCE_PRESS       = 0x09 // 按下亮 (press to light)
	SEQUENCE_SNAKE       = 0x0A // 往中间蛇形跑灯 (center snake)
	SEQUENCE_FOUNTAIN    = 0x0B // 七彩喷泉 (color fountain)
	SEQUENCE_LASER       = 0x0C // 按键激光 (key laser)
	SEQUENCE_FISH        = 0x0D // 发光的鱼 (glowing fish)
	SEQUENCE_SURF_CROSS  = 0x0E // 交叉冲浪 (cross surfing)
	SEQUENCE_HEART       = 0x0F // 爱心 (heart)
	SEQUENCE_TRAFFIC     = 0x10 // 交通 (traffic)
	SEQUENCE_SNAKE_OLD   = 0x11 // 贪吃蛇 (snake)
	SEQUENCE_RAINDROP    = 0x12 // 雨滴 (raindrop)
	SEQUENCE_CUSTOM      = 0x13 // custom colors (turbo or standard)
)

// Color index values for G60
const (
	COLOR_RED    = 0x01
	COLOR_ORANGE = 0x02
	COLOR_YELLOW = 0x03
	COLOR_GREEN  = 0x04
	COLOR_CYAN   = 0x05
	COLOR_BLUE   = 0x06
	COLOR_PURPLE = 0x07
	COLOR_WHITE  = 0x08
)

const LAYOUT_SIZE = 126

type Layout struct {
	names [LAYOUT_SIZE]string
	index map[string]int
}

func newLayout(names [LAYOUT_SIZE]string) Layout {
	l := Layout{names: names, index: make(map[string]int, LAYOUT_SIZE)}
	for i, n := range names {
		if n != "" {
			l.index[n] = i
		}
	}
	return l
}

func (l Layout) IndexOf(name string) int {
	if idx, ok := l.index[name]; ok {
		return idx
	}
	return -1
}

var a75Layout = newLayout([LAYOUT_SIZE]string{
	0: "ESC",
	2: "F1", 3: "F2", 4: "F3", 5: "F4", 6: "F5", 7: "F6", 8: "F7", 9: "F8", 10: "F9", 11: "F10", 12: "F11", 13: "F12",
	14: "DEL",
	21: "TILDE",
	22: "1", 23: "2", 24: "3", 25: "4", 26: "5", 27: "6", 28: "7", 29: "8", 30: "9", 31: "0",
	32: "MINUS", 33: "PLUS", 34: "BACK",
	36: "HOME",
	42: "TAB",
	43: "Q", 44: "W", 45: "E", 46: "R", 47: "T", 48: "Y", 49: "U", 50: "I", 51: "O", 52: "P",
	53: "BRKTS_L", 54: "BRKTS_R", 55: "SLASH_K29",
	57: "PGUP",
	63: "CAPS",
	64: "A", 65: "S", 66: "D", 67: "F", 68: "G", 69: "H", 70: "J", 71: "K", 72: "L",
	73: "COLON", 74: "QOTATN", 76: "RETURN",
	78: "PGDN",
	84: "SHF_L",
	86: "Z", 87: "X", 88: "C", 89: "V", 90: "B", 91: "N", 92: "M",
	93: "COMMA", 94: "PERIOD", 95: "SLASH",
	97: "SHF_R", 98: "ARR_UP", 99: "END",
	105: "CTRL_L", 106: "WIN_L", 107: "ALT_L",
	111: "SPACE",
	115: "ALT_R", 116: "FN1", 117: "MENU",
	119: "ARR_L", 120: "ARR_DW", 121: "ARR_R",
})

var g60Layout = newLayout([LAYOUT_SIZE]string{
	21: "ESC",
	22: "1", 23: "2", 24: "3", 25: "4", 26: "5", 27: "6", 28: "7", 29: "8", 30: "9", 31: "0",
	32: "MINUS", 33: "PLUS", 34: "BACK",
	42: "TAB",
	43: "Q", 44: "W", 45: "E", 46: "R", 47: "T", 48: "Y",
	49: "U", 50: "I", 51: "O", 52: "P",
	53: "BRKTS_L", 54: "BRKTS_R", 55: "SLASH_K29",
	63: "CAPS",
	64: "A", 65: "S", 66: "D", 67: "F", 68: "G", 69: "H",
	70: "J", 71: "K", 72: "L",
	73: "COLON", 74: "QOTATN", 76: "RETURN",
	84: "SHF_L",
	86: "Z", 87: "X", 88: "C", 89: "V",
	90: "B", 91: "N", 92: "M",
	93: "COMMA", 94: "PERIOD", 95: "SLASH",
	97: "SHF_R",
	105: "CTRL_L", 106: "WIN_L", 107: "ALT_L",
	111: "SPACE",
	115: "ALT_R", 116: "FN1", 117: "FN2", 118: "CTRL_R",
})

var g65Layout = newLayout([LAYOUT_SIZE]string{
	21: "ESC",
	22: "1", 23: "2", 24: "3", 25: "4", 26: "5", 27: "6", 28: "7", 29: "8", 30: "9", 31: "0",
	32: "MINUS", 33: "PLUS", 34: "BACK", 35: "DEL",
	42: "TAB",
	43: "Q", 44: "W", 45: "E", 46: "R", 47: "T", 48: "Y",
	49: "U", 50: "I", 51: "O", 52: "P",
	53: "BRKTS_L", 55: "SLASH_K29", 56: "END",
	63: "CAPS",
	64: "A", 65: "S", 66: "D", 67: "F", 68: "G", 69: "H",
	70: "J", 71: "K", 72: "L",
	73: "COLON", 74: "QOTATN", 76: "RETURN", 77: "PGUP",
	84: "SHF_L",
	86: "Z", 87: "X", 88: "C", 89: "V",
	90: "B", 91: "N", 92: "M",
	93: "COMMA", 94: "PERIOD", 95: "SLASH",
	96: "SHF_R",
	97: "ARR_UP",
	98: "PGDN",
	105: "CTRL_L", 106: "WIN_L", 107: "ALT_L",
	111: "SPACE",
	114: "ALT_R", 115: "FN1", 116: "MENU",
	117: "ARR_L", 118: "ARR_DW", 119: "ARR_R",
})

var g75Layout = newLayout([LAYOUT_SIZE]string{
	0: "ESC",
	1: "F1", 2: "F2", 3: "F3", 4: "F4", 5: "F5", 6: "F6",
	7: "F7", 8: "F8", 9: "F9", 10: "F10", 11: "F11", 12: "F12",
	13: "PRINT", 14: "INS", 15: "DEL",
	21: "TILDE",
	22: "1", 23: "2", 24: "3", 25: "4", 26: "5", 27: "6",
	28: "7", 29: "8", 30: "9", 31: "0", 32: "MINUS", 33: "PLUS", 34: "BACK",
	36: "HOME",
	42: "TAB",
	43: "Q", 44: "W", 45: "E", 46: "R", 47: "T", 48: "Y",
	49: "U", 50: "I", 51: "O", 52: "P",
	53: "BRKTS_L", 54: "BRKTS_R", 55: "SLASH_K29",
	57: "PGUP",
	63: "CAPS",
	64: "A", 65: "S", 66: "D", 67: "F", 68: "G", 69: "H",
	70: "J", 71: "K", 72: "L",
	73: "COLON", 74: "QOTATN", 76: "RETURN",
	78: "PGDN",
	84: "SHF_L",
	86: "Z", 87: "X", 88: "C", 89: "V", 90: "B", 91: "N", 92: "M",
	93: "COMMA", 94: "PERIOD", 95: "SLASH",
	96: "SHF_R",
	97: "ARR_UP",
	99: "END",
	105: "CTRL_L", 106: "WIN_L", 107: "ALT_L",
	111: "SPACE",
	114: "ALT_R", 115: "FN1",
	117: "MENU",
	118: "ARR_L", 119: "ARR_DW", 120: "ARR_R",
})

func GetLayout(model string) Layout {
	switch model {
	case KEYBOARD_G60:
		return g60Layout
	case KEYBOARD_G65:
		return g65Layout
	case KEYBOARD_G75:
		return g75Layout
	default:
		return a75Layout
	}
}

// G60 LED grid in row-major order, matching the webdriver's getG60() iteration
var G60_LED_GRID = [][]int{
	{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34},
	{42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55},
	{63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 76},
	{84, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 97},
	{105, 106, 107, 111, 115, 116, 117, 118},
}

var WASD_KEYS = []int{44, 64, 65, 66}
var NUMERALS_KEYS = []int{22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
var CHARACTER_KEYS = []int{
	43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
	64, 65, 66, 67, 68, 69, 70, 71, 72,
	86, 87, 88, 89, 90, 91, 92,
}
