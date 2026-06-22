package driver

import (
	"bytes"
	"strconv"
)

func BoolToByte(b bool) byte {
	if b {
		return 0x01
	}
	return 0x00
}

func DetectKeyboardModel(modelBytes []byte) (string, int) { // Model and type
	A75 := [][]byte{
		{0x0b, 0x01, 0x01},
		{0x0b, 0x04, 0x01},
	}

	A75Pro := [][]byte{
		{0x0b, 0x04, 0x03},
	}

	G65 := [][]byte{
		{0x0f, 0x01, 0x01},
		{0x0b, 0x02, 0x01},
	}

	G60 := [][]byte{
		{0x0b, 0x03, 0x01},
	}

	G75 := [][]byte{
		{0x0b, 0x04, 0x05},
	}

	for _, modelA75 := range A75 {
		if bytes.Equal(modelA75, modelBytes) {
			return KEYBOARD_A75, 75
		}
	}

	for _, modelA75Pro := range A75Pro {
		if bytes.Equal(modelA75Pro, modelBytes) {
			return KEYBOARD_A75PRO, 750
		}
	}

	for _, modelG75 := range G75 {
		if bytes.Equal(modelG75, modelBytes) {
			return KEYBOARD_G75, 754
		}
	}

	for _, modelG65 := range G65 {
		if bytes.Equal(modelG65, modelBytes) {
			return KEYBOARD_G65, 65
		}
	}

	for _, modelG60 := range G60 {
		if bytes.Equal(modelG60, modelBytes) {
			return KEYBOARD_G60, 60
		}
	}

	// Default to A75 if no match found
	return KEYBOARD_A75, 75
}

func GetIndexByKey(key string, model string) int {
	return GetLayout(model).IndexOf(key)
}

func GetKeyByIndex(index int, model string) string {
	if index >= 0 && index < LAYOUT_SIZE {
		return GetLayout(model).names[index]
	}
	return ""
}

// ResolveKeyToIndex returns the firmware index for a key.
// Tries name lookup first, then parses as a numeric firmware index.
func ResolveKeyToIndex(key string, model string) (int, bool) {
	idx := GetIndexByKey(key, model)
	if idx >= 0 {
		return idx, true
	}
	if n, err := strconv.Atoi(key); err == nil {
		if n >= 0 && n < LAYOUT_SIZE {
			return n, true
		}
	}
	return -1, false
}

func GetLayoutKeys(model string) []string {
	l := GetLayout(model)
	keys := make([]string, 0, LAYOUT_SIZE)
	for _, name := range l.names {
		if name != "" {
			keys = append(keys, name)
		}
	}
	return keys
}

func GetRowByIndex(index int) int {
	return index / KEYS_PER_ROW
}

func ActuationFloatToByte(actuation float32) byte {
	var normalized int = int(actuation * 10)
	if normalized < 0 {
		normalized = 1
	} // maniacs

	if normalized > 255 {
		normalized = 255
	} // WHY WOULD ANYONE EVEN TRY TO SET THIS

	return byte(normalized)
}
