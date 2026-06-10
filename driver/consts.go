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

// Imagine this is a const
var KEYBOARD_LAYOUT = []string{
	"ESC", "", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12", "KP7", "KP8", "KP9", "", "", "", "",
	"TILDE", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "MINUS", "PLUS", "BACK", "KP4", "KP5", "KP6", "", "", "", "",
	"TAB", "Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "BRKTS_L", "BRKTS_R", "SLASH_K29", "KP1", "KP2", "KP3", "", "", "", "",
	"CAPS", "A", "S", "D", "F", "G", "H", "J", "K", "L", "COLON", "QOTATN", "", "RETURN", "", "KP0", "KP_DEL", "", "", "", "",
	"SHF_L", "EUR_K45", "Z", "X", "C", "V", "B", "N", "M", "COMMA", "PERIOD", "SLASH", "", "SHF_R", "ARR_UP", "", "NUMS", "", "", "", "",
	"CTRL_L", "WIN_L", "ALT_L", "", "", "", "SPACE", "", "", "", "ALT_R", "FN1", "APP", "ARR_L", "ARR_DW", "ARR_R", "CTRL_R", "", "", "", "",
}

// G60 LED grid in row-major order, matching the webdriver's getG60() iteration
var G60_LED_GRID = [][]int{
	{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34},
	{42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55},
	{63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 76},
	{84, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 97},
	{105, 106, 107, 111, 115, 116, 117, 118},
}

// G60 key name to firmware index. KEYBOARD_LAYOUT is a generic full-size layout
// where "ESC" is at index 0, but the G60 firmware has ESC at index 21.
var G60_KEY_INDEX = map[string]int{
	"ESC": 21, "1": 22, "2": 23, "3": 24, "4": 25, "5": 26, "6": 27,
	"7": 28, "8": 29, "9": 30, "0": 31,
	"MINUS": 32, "PLUS": 33, "BACK": 34,
	"TAB": 42,
	"Q": 43, "W": 44, "E": 45, "R": 46, "T": 47, "Y": 48,
	"U": 49, "I": 50, "O": 51, "P": 52,
	"BRKTS_L": 53, "BRKTS_R": 54, "SLASH_K29": 55,
	"CAPS": 63,
	"A": 64, "S": 65, "D": 66, "F": 67, "G": 68, "H": 69,
	"J": 70, "K": 71, "L": 72,
	"COLON": 73, "QOTATN": 74, "RETURN": 76,
	"SHF_L": 84,
	"Z": 86, "X": 87, "C": 88, "V": 89,
	"B": 90, "N": 91, "M": 92,
	"COMMA": 93, "PERIOD": 94, "SLASH": 95,
	"SHF_R": 97,
	"CTRL_L": 105, "WIN_L": 106, "ALT_L": 107,
	"SPACE": 111,
	"ALT_R": 115, "FN1": 116, "FN2": 117, "CTRL_R": 118,
}

var WASD_KEYS = []int{44, 64, 65, 66}
var NUMERALS_KEYS = []int{22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
var CHARACTER_KEYS = []int{
	43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
	64, 65, 66, 67, 68, 69, 70, 71, 72,
	86, 87, 88, 89, 90, 91, 92,
}
