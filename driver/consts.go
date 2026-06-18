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

type Layout [LAYOUT_SIZE]string

func (l Layout) IndexOf(name string) int {
	for i, n := range l {
		if n == name {
			return i
		}
	}
	return -1
}

var a75Layout = Layout{
	0: "ESC",
	2: "F1", 3: "F2", 4: "F3", 5: "F4", 6: "F5", 7: "F6", 8: "F7", 9: "F8", 10: "F9", 11: "F10", 12: "F11", 13: "F12",
	14: "KP7", 15: "KP8", 16: "KP9",
	21: "TILDE",
	22: "1", 23: "2", 24: "3", 25: "4", 26: "5", 27: "6", 28: "7", 29: "8", 30: "9", 31: "0",
	32: "MINUS", 33: "PLUS", 34: "BACK",
	35: "KP4", 36: "KP5", 37: "KP6",
	42: "TAB",
	43: "Q", 44: "W", 45: "E", 46: "R", 47: "T", 48: "Y", 49: "U", 50: "I", 51: "O", 52: "P",
	53: "BRKTS_L", 54: "BRKTS_R", 55: "SLASH_K29",
	56: "KP1", 57: "KP2", 58: "KP3",
	63: "CAPS",
	64: "A", 65: "S", 66: "D", 67: "F", 68: "G", 69: "H", 70: "J", 71: "K", 72: "L",
	73: "COLON", 74: "QOTATN", 76: "RETURN",
	77: "KP0", 78: "KP_DEL",
	84: "SHF_L",
	85: "EUR_K45",
	86: "Z", 87: "X", 88: "C", 89: "V", 90: "B", 91: "N", 92: "M",
	93: "COMMA", 94: "PERIOD", 95: "SLASH",
	97: "SHF_R", 98: "ARR_UP", 100: "NUMS",
	105: "CTRL_L", 106: "WIN_L", 107: "ALT_L",
	111: "SPACE",
	115: "ALT_R", 116: "FN1", 117: "APP", 118: "ARR_L", 119: "ARR_DW", 120: "ARR_R", 121: "CTRL_R",
}

var g60Layout = Layout{
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
}

func GetLayout(model string) Layout {
	switch model {
	case KEYBOARD_G60:
		return g60Layout
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
