package driver

type RemapKey struct {
	KeyCmd  byte
	KeyCode byte
	KeyType byte
}

var actionLookup = map[string]RemapKey{
	"CTRL_L":  {0xFC, 1, 1},
	"SHIFT_L": {0xFC, 2, 1},
	"SHF_L":   {0xFC, 2, 1},
	"ALT_L":   {0xFC, 4, 1},
	"WIN_L":   {0xFC, 227, 0},
	"CTRL_R":  {0xFC, 16, 1},
	"SHIFT_R": {0xFC, 32, 1},
	"SHF_R":   {0xFC, 32, 1},
	"ALT_R":   {0xFC, 64, 1},
	"FN1":     {0xFF, 0x80, 1},
	"FN2":     {0xFE, 0x80, 1},

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
	"SLASH":    {0xFC, 0x38, 0},
	"SLASH_K29": {0xFC, 0x31, 0},
	"EUR_K45":  {0xFC, 0x64, 0},
	"CAPS":     {0xFC, 0x39, 0},

	"F1":  {0xFC, 0x3A, 0}, "F2": {0xFC, 0x3B, 0},
	"F3":  {0xFC, 0x3C, 0}, "F4": {0xFC, 0x3D, 0},
	"F5":  {0xFC, 0x3E, 0}, "F6": {0xFC, 0x3F, 0},
	"F7":  {0xFC, 0x40, 0}, "F8": {0xFC, 0x41, 0},
	"F9":  {0xFC, 0x42, 0}, "F10": {0xFC, 0x43, 0},
	"F11": {0xFC, 0x44, 0}, "F12": {0xFC, 0x45, 0},

	"PRINT":  {0xFC, 0x46, 0},
	"SCRLK":  {0xFC, 0x47, 0},
	"PAUSE":  {0xFC, 0x48, 0},
	"INS":    {0xFC, 0x49, 0},
	"HOME":   {0xFC, 0x4A, 0},
	"PGUP":   {0xFC, 0x4B, 0},
	"DEL":    {0xFC, 0x4C, 0},
	"END":    {0xFC, 0x4D, 0},
	"PGDN":   {0xFC, 0x4E, 0},
	"ARR_R":  {0xFC, 0x4F, 0},
	"ARR_L":  {0xFC, 0x50, 0},
	"ARR_DW": {0xFC, 0x51, 0},
	"ARR_UP": {0xFC, 0x52, 0},
	"APP":    {0xFC, 0x65, 0},

	"NUMS":   {0xFC, 0x53, 0},
	"NUM7":    {0xFC, 0x5F, 0},
	"NUM8":    {0xFC, 0x60, 0},
	"NUM9":    {0xFC, 0x61, 0},
	"NUM4":    {0xFC, 0x5C, 0},
	"NUM5":    {0xFC, 0x5D, 0},
	"NUM6":    {0xFC, 0x5E, 0},
	"NUM1":    {0xFC, 0x59, 0},
	"NUM2":    {0xFC, 0x5A, 0},
	"NUM3":    {0xFC, 0x5B, 0},
	"NUM0":    {0xFC, 0x62, 0},
	"KP_DEL":   {0xFC, 0x63, 0},
	"KP_ENTER": {0xFC, 0x58, 0},
	"KP_PLUS":  {0xFC, 0x57, 0},
	"KP_MINUS": {0xFC, 0x56, 0},
	"KP_MULT":  {0xFC, 0x55, 0},
	"KP_DIV":   {0xFC, 0x54, 0},

	"K45":   {0xFC, 0x64, 0},
	"K56":   {0xFC, 0x87, 0},
	"LOOP":  {0xFC, 0x8A, 0},
	"NOLOOP": {0xFC, 0x8B, 0},
	"KANA":  {0xFC, 0x88, 0},

	"VOL_UP":   {0x90, 0, 3},
	"VOL_DN":   {0x91, 0, 3},
	"MUTE":     {0x92, 0, 3},
	"PLAY":     {0x93, 0, 3},
	"STOP":     {0x94, 0, 3},
	"PREV":     {0x95, 0, 3},
	"NEXT":     {0x96, 0, 3},
	"MS_L":     {0xB0, 0, 3},
	"MS_R":     {0xB1, 0, 3},
	"MS_M":     {0xB2, 0, 3},
	"MS_SCR_U": {0xB5, 0, 3},
	"MS_SCR_D": {0xB6, 0, 3},
	"MS_SCR_L": {0xB7, 0, 3},
	"MS_SCR_R": {0xB8, 0, 3},

	"WIN_LOCK":    {0xFB, 0x10, 2},
	"RDT":         {0xFB, 0x4E, 2},
	"LW":          {0xFB, 0x4F, 2},
	"RT":          {0xFD, 0x20, 2},
	"RT_MATCH":    {0xFD, 0x0E, 2},
	"RT_EXTREME":  {0xFD, 0x30, 2},
	"RT_STANDARD": {0xFD, 0x32, 2},
	"RT_COMPETE":  {0xFD, 0x31, 2},
	"LIGHT_SW":    {0xFB, 0x0E, 2},
	"LIGHT_CYCLE": {0xFB, 0x00, 2},
	"LIGHT_NEXT":  {0xFB, 0x01, 2},
	"LIGHT_PREV":  {0xFB, 0x02, 2},
	"LIGHT_BR_INC": {0xFB, 0x04, 2},
	"LIGHT_BR_DEC": {0xFB, 0x05, 2},
	"LIGHT_SP_INC": {0xFB, 0x07, 2},
	"LIGHT_SP_DEC": {0xFB, 0x08, 2},
	"LIGHT_COLOR": {0xFB, 0x09, 2},
	"LIGHT_PAUSE": {0xFB, 0x0D, 2},
	"LIGHT_CTRL":  {0xFB, 0x30, 2},
	"COLOR_CTRL":  {0xFB, 0x31, 2},
	"SPEED_CTRL":  {0xFB, 0x32, 2},
	"BR_CTRL":     {0xFB, 0x33, 2},
}

func GetRemapAction(name string) (RemapKey, bool) {
	k, ok := actionLookup[name]
	return k, ok
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
	"ARR_LEFT": "\U0000F060",
	"ARR_DOWN": "\U0000F063",
	"ARR_RIGHT":"\U0000F061",
}

func ActionDisplayName(action string) string {
	if d, ok := actionDisplayNames[action]; ok {
		return d
	}
	return action
}

func GetRemapIndexByKey(key string, model string) int {
	return GetLayout(model).IndexOf(key)
}

func BuildRemapChunk(layer, chunkNum byte, keys map[int]*RemapKey) []byte {
	report := make([]byte, 63)
	report[0] = PACKET_IDENTITY
	report[1] = 0x02
	report[2] = 0x04
	report[3] = chunkNum
	report[4] = 14
	report[5] = layer

	start := int(chunkNum-1) * 9
	r := 6
	for w := 0; w < 9; w++ {
		s := start + w
		if s >= 126 {
			break
		}
		k, ok := keys[s]
		if ok && k != nil {
			report[r] = byte(s)
			r++
			switch k.KeyType {
			case 0:
				report[r] = k.KeyCmd
				r += 2
				report[r] = k.KeyCode
				r += 3
			case 1, 2:
				report[r] = k.KeyCmd
				r++
				report[r] = k.KeyCode
				r += 4
			case 3:
				report[r] = k.KeyCmd
				r += 5
			case 4:
				report[r] = 0xF8
				r++
				report[r] = k.KeyCode
				r++
				report[r] = 0
				r++
				report[r] = 0
				r++
				report[r] = 0
				r++
			default:
				report[r] = k.KeyCmd
				r += 5
			}
		} else {
			r += 6
		}
	}
	report[60] = 0xA5

	return report
}
