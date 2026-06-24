package driver

type RemapKey struct {
	KeyCmd  byte
	KeyCode byte
	KeyType byte
}

// Action name constants for keys not physically present on G60 + all lighting/feature actions.
// Values are the canonical actionLookup keys (backward-compatible with profile JSON).
const (
	KeyCtrlL = "CTRL_L"
	KeyShiftL = "SHIFT_L"
	KeyAltL   = "ALT_L"
	KeyWinL   = "WIN_L"
	KeyCtrlR  = "CTRL_R"
	KeyShiftR = "SHIFT_R"
	KeyAltR   = "ALT_R"
	KeyFn1    = "FN1"
	KeyFn2    = "FN2"

	KeyF1  = "F1"
	KeyF2  = "F2"
	KeyF3  = "F3"
	KeyF4  = "F4"
	KeyF5  = "F5"
	KeyF6  = "F6"
	KeyF7  = "F7"
	KeyF8  = "F8"
	KeyF9  = "F9"
	KeyF10 = "F10"
	KeyF11 = "F11"
	KeyF12 = "F12"

	KeyPrint = "PRINT"
	KeyScrlk = "SCRLK"
	KeyPause = "PAUSE"
	KeyIns   = "INS"
	KeyHome  = "HOME"
	KeyPgUp  = "PGUP"
	KeyDel   = "DEL"
	KeyEnd   = "END"
	KeyPgDn  = "PGDN"

	KeyArrRight = "ARR_R"
	KeyArrLeft  = "ARR_L"
	KeyArrDown  = "ARR_DW"
	KeyArrUp    = "ARR_UP"

	KeyNums    = "NUMS"
	KeyNum7    = "NUM7"
	KeyNum8    = "NUM8"
	KeyNum9    = "NUM9"
	KeyNum4    = "NUM4"
	KeyNum5    = "NUM5"
	KeyNum6    = "NUM6"
	KeyNum1    = "NUM1"
	KeyNum2    = "NUM2"
	KeyNum3    = "NUM3"
	KeyNum0    = "NUM0"
	KeyKpDel   = "KP_DEL"
	KeyKpEnter = "KP_ENTER"
	KeyKpPlus  = "KP_PLUS"
	KeyKpMinus = "KP_MINUS"
	KeyKpMult  = "KP_MULT"
	KeyKpDiv   = "KP_DIV"

	KeyVolUp  = "VOL_UP"
	KeyVolDn  = "VOL_DN"
	KeyMute   = "MUTE"
	KeyPlay   = "PLAY"
	KeyStop   = "STOP"
	KeyPrev   = "PREV"
	KeyNext   = "NEXT"

	KeyMsL    = "MS_L"
	KeyMsR    = "MS_R"
	KeyMsM    = "MS_M"
	KeyMsScrU = "MS_SCR_U"
	KeyMsScrD = "MS_SCR_D"
	KeyMsScrL = "MS_SCR_L"
	KeyMsScrR = "MS_SCR_R"

	KeyApp       = "APP"
	KeyK45       = "K45"
	KeyK56       = "K56"
	KeyLoop      = "LOOP"
	KeyNoLoop    = "NOLOOP"
	KeyKana      = "KANA"
	KeySlashK29  = "SLASH_K29"
	KeyEurK45    = "EUR_K45"

	KeyLightModeSw    = "LIGHT_SW"
	KeyLightModeCycle = "LIGHT_CYCLE"
	KeyLightModeNext  = "LIGHT_NEXT"
	KeyLightModePrev  = "LIGHT_PREV"
	KeyLightLumiInc   = "LIGHT_BR_INC"
	KeyLightLumiDec   = "LIGHT_BR_DEC"
	KeyLightSpeedInc  = "LIGHT_SP_INC"
	KeyLightSpeedDec  = "LIGHT_SP_DEC"
	KeyLightColorCycle = "LIGHT_COLOR"
	KeyLightPause     = "LIGHT_PAUSE"
	KeyLightCtrl      = "LIGHT_CTRL"
	KeyColorCtrl      = "COLOR_CTRL"
	KeySpeedCtrl      = "SPEED_CTRL"
	KeyBrCtrl         = "BR_CTRL"
	KeyWinLock        = "WIN_LOCK"

	KeyRdt        = "RDT"
	KeyLw         = "LW"
	KeyTurbo      = "Turbo"
	KeyRt         = "RT"
	KeyRtMatch    = "RT_MATCH"
	KeyRtExtreme  = "RT_EXTREME"
	KeyRtStandard = "RT_STANDARD"
	KeyRtCompete  = "RT_COMPETE"
)

var actionLookup = map[string]RemapKey{
	KeyCtrlL:   {0xFC, 1, 1},
	KeyShiftL:  {0xFC, 2, 1},
	KeyAltL:    {0xFC, 4, 1},
	KeyWinL:    {0xFC, 227, 0},
	KeyCtrlR:   {0xFC, 16, 1},
	KeyShiftR:  {0xFC, 32, 1},
	KeyAltR:    {0xFC, 64, 1},
	KeyFn1:     {0xFF, 0x80, 1},
	KeyFn2:     {0xFE, 0x80, 1},

	"A": {0xFC, 0x04, 0}, "B": {0xFC, 0x05, 0},
	"C": {0xFC, 0x06, 0}, "D": {0xFC, 0x07, 0},
	"E": {0xFC, 0x08, 0}, "F": {0xFC, 0x09, 0},
	"G": {0xFC, 0x0A, 0}, "H": {0xFC, 0x0B, 0},
	"I": {0xFC, 0x0C, 0}, "J": {0xFC, 0x0D, 0},
	"K": {0xFC, 0x0E, 0}, "L": {0xFC, 0x0F, 0},
	"M": {0xFC, 0x10, 0}, "N": {0xFC, 0x11, 0},
	"O": {0xFC, 0x12, 0}, "P": {0xFC, 0x13, 0},
	"Q": {0xFC, 0x14, 0}, "R": {0xFC, 0x15, 0},
	"S": {0xFC, 0x16, 0}, "T": {0xFC, 0x17, 0},
	"U": {0xFC, 0x18, 0}, "V": {0xFC, 0x19, 0},
	"W": {0xFC, 0x1A, 0}, "X": {0xFC, 0x1B, 0},
	"Y": {0xFC, 0x1C, 0}, "Z": {0xFC, 0x1D, 0},

	"1": {0xFC, 0x1E, 0}, "2": {0xFC, 0x1F, 0},
	"3": {0xFC, 0x20, 0}, "4": {0xFC, 0x21, 0},
	"5": {0xFC, 0x22, 0}, "6": {0xFC, 0x23, 0},
	"7": {0xFC, 0x24, 0}, "8": {0xFC, 0x25, 0},
	"9": {0xFC, 0x26, 0}, "0": {0xFC, 0x27, 0},

	"ENTER":   {0xFC, 0x28, 0},
	"ESC":     {0xFC, 0x29, 0},
	"BACK":    {0xFC, 0x2A, 0},
	"TAB":     {0xFC, 0x2B, 0},
	"SPACE":   {0xFC, 0x2C, 0},
	"MINUS":   {0xFC, 0x2D, 0},
	"PLUS":    {0xFC, 0x2E, 0},
	"BRKTS_L": {0xFC, 0x2F, 0},
	"BRKTS_R": {0xFC, 0x30, 0},
	"COLON":   {0xFC, 0x33, 0},
	"QOTATN":  {0xFC, 0x34, 0},
	"TILDE":   {0xFC, 0x35, 0},
	"COMMA":   {0xFC, 0x36, 0},
	"PERIOD":  {0xFC, 0x37, 0},
	"SLASH":   {0xFC, 0x38, 0},
	KeySlashK29: {0xFC, 0x31, 0},
	KeyEurK45:   {0xFC, 0x64, 0},
	"CAPS":    {0xFC, 0x39, 0},

	KeyF1: {0xFC, 0x3A, 0}, KeyF2: {0xFC, 0x3B, 0},
	KeyF3: {0xFC, 0x3C, 0}, KeyF4: {0xFC, 0x3D, 0},
	KeyF5: {0xFC, 0x3E, 0}, KeyF6: {0xFC, 0x3F, 0},
	KeyF7: {0xFC, 0x40, 0}, KeyF8: {0xFC, 0x41, 0},
	KeyF9: {0xFC, 0x42, 0}, KeyF10: {0xFC, 0x43, 0},
	KeyF11: {0xFC, 0x44, 0}, KeyF12: {0xFC, 0x45, 0},

	KeyPrint: {0xFC, 0x46, 0},
	KeyScrlk: {0xFC, 0x47, 0},
	KeyPause: {0xFC, 0x48, 0},
	KeyIns:   {0xFC, 0x49, 0},
	KeyHome:  {0xFC, 0x4A, 0},
	KeyPgUp:  {0xFC, 0x4B, 0},
	KeyDel:   {0xFC, 0x4C, 0},
	KeyEnd:   {0xFC, 0x4D, 0},
	KeyPgDn:  {0xFC, 0x4E, 0},
	KeyArrRight: {0xFC, 0x4F, 0},
	KeyArrLeft: {0xFC, 0x50, 0},
	KeyArrDown: {0xFC, 0x51, 0},
	KeyArrUp:  {0xFC, 0x52, 0},
	KeyApp:    {0xFC, 0x65, 0},

	KeyNums:    {0xFC, 0x53, 0},
	KeyNum7:    {0xFC, 0x5F, 0},
	KeyNum8:    {0xFC, 0x60, 0},
	KeyNum9:    {0xFC, 0x61, 0},
	KeyNum4:    {0xFC, 0x5C, 0},
	KeyNum5:    {0xFC, 0x5D, 0},
	KeyNum6:    {0xFC, 0x5E, 0},
	KeyNum1:    {0xFC, 0x59, 0},
	KeyNum2:    {0xFC, 0x5A, 0},
	KeyNum3:    {0xFC, 0x5B, 0},
	KeyNum0:    {0xFC, 0x62, 0},
	KeyKpDel:   {0xFC, 0x63, 0},
	KeyKpEnter: {0xFC, 0x58, 0},
	KeyKpPlus:  {0xFC, 0x57, 0},
	KeyKpMinus: {0xFC, 0x56, 0},
	KeyKpMult:  {0xFC, 0x55, 0},
	KeyKpDiv:   {0xFC, 0x54, 0},

	KeyK45:    {0xFC, 0x64, 0},
	KeyK56:    {0xFC, 0x87, 0},
	KeyLoop:   {0xFC, 0x8A, 0},
	KeyNoLoop: {0xFC, 0x8B, 0},
	KeyKana:   {0xFC, 0x88, 0},

	KeyVolUp:  {0x90, 0, 3},
	KeyVolDn:  {0x91, 0, 3},
	KeyMute:   {0x92, 0, 3},
	KeyPlay:   {0x93, 0, 3},
	KeyStop:   {0x94, 0, 3},
	KeyPrev:   {0x95, 0, 3},
	KeyNext:   {0x96, 0, 3},
	KeyMsL:    {0xB0, 0, 3},
	KeyMsR:    {0xB1, 0, 3},
	KeyMsM:    {0xB2, 0, 3},
	KeyMsScrU: {0xB5, 0, 3},
	KeyMsScrD: {0xB6, 0, 3},
	KeyMsScrL: {0xB7, 0, 3},
	KeyMsScrR: {0xB8, 0, 3},

	KeyWinLock: {0xFB, 0x10, 2},
	KeyRdt:     {0xFB, 0x4E, 2},
	KeyLw:      {0xFB, 0x4F, 2},
	KeyRt:      {0xFD, 0x20, 2},
	KeyRtMatch: {0xFD, 0x0E, 2},
	KeyRtExtreme:  {0xFD, 0x30, 2},
	KeyRtStandard: {0xFD, 0x32, 2},
	KeyRtCompete:  {0xFD, 0x31, 2},
	KeyTurbo:  {0xFD, 0x08, 2},

	KeyLightModeSw:    {0xFB, 0x0E, 2},
	KeyLightModeCycle: {0xFB, 0x00, 2},
	KeyLightModeNext:  {0xFB, 0x01, 2},
	KeyLightModePrev:  {0xFB, 0x02, 2},
	KeyLightLumiInc:   {0xFB, 0x04, 2},
	KeyLightLumiDec:   {0xFB, 0x05, 2},
	KeyLightSpeedInc:  {0xFB, 0x07, 2},
	KeyLightSpeedDec:  {0xFB, 0x08, 2},
	KeyLightColorCycle: {0xFB, 0x09, 2},
	KeyLightPause:     {0xFB, 0x0D, 2},
	KeyLightCtrl:      {0xFB, 0x30, 2},
	KeyColorCtrl:      {0xFB, 0x31, 2},
	KeySpeedCtrl:      {0xFB, 0x32, 2},
	KeyBrCtrl:         {0xFB, 0x33, 2},
}

var actionDisplayNames = map[string]string{
	"VOL_UP":   "\U000F075D",
	"VOL_DN":   "\U000F075E",
	"PREV":     "\U000F04AE",
	"NEXT":     "\U000F04AD",
	"PLAY":     "\U000F040A",
	"STOP":     "\U000F03E4",
	"MUTE":     "\U000F075F",
	"ARR_UP":   "\U0000F062",
	"ARR_L": "\U0000F060",
	"ARR_DW": "\U0000F063",
	"ARR_R": "\U0000F061",

	KeyLightModeSw:     "Light\nToggle",
	KeyTurbo:           "Turbo\nToggle",
	KeyLightColorCycle: "Color\nCycle",
	KeyLightModeCycle:  "Light\nCycle",
	KeyLightLumiInc:    "Light\n+",
	KeyLightLumiDec:    "Light\n-",
	KeyLightSpeedInc:   "Speed\n+",
	KeyLightSpeedDec:   "Speed\n-",
	KeyLightModeNext:   "Light\nNext",
	KeyLightModePrev:   "Light\nPrev",
	KeyLightPause:      "Light\nPause",
	KeyWinLock:         "Win\nLock",
	KeyRdt:             "RDT",
	KeyLw:              "LW",
	KeyFn1:             "Fn",
	KeyFn2:             "Menu",
	KeyPrint:           "Print\nScr",
	KeyScrlk:           "Scroll\nLock",
	KeyPause:           "Pause",
	KeyIns:             "Insert",
	KeyHome:            "Home",
	KeyPgUp:            "Page\nUp",
	KeyDel:             "Delete",
	KeyEnd:             "End",
	KeyPgDn:            "Page\nDown",
}

func ActionDisplayName(action string) string {
	if d, ok := actionDisplayNames[action]; ok {
		return d
	}
	return action
}

func GetRemapAction(name string) (RemapKey, bool) {
	k, ok := actionLookup[name]
	return k, ok
}

func GetRemapIndexByKey(key string, model string) int {
	return GetLayout(model).IndexOf(key)
}
