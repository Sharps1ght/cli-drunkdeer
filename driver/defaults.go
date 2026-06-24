package driver

var defaultFnKeys = map[string]map[int]bool{
	KEYBOARD_G60: {
		21: true, 22: true, 23: true, 24: true, 25: true,
		26: true, 27: true, 28: true, 29: true, 30: true,
		31: true, 32: true, 33: true,
		44: true,
		52: true, 53: true, 54: true,
		64: true, 65: true, 66: true,
		72: true, 73: true, 74: true,
		86: true, 87: true, 88: true, 89: true, 90: true, 91: true, 92: true,
		93: true, 94: true, 95: true,
	},
	KEYBOARD_G65: {
		21: true, 22: true, 23: true, 24: true, 25: true,
		26: true, 27: true, 28: true, 29: true, 30: true,
		31: true, 32: true, 33: true,
		35: true,
		54: true, 56: true,
		77: true,
		86: true, 87: true, 88: true, 89: true, 90: true, 91: true, 92: true,
		98: true,
	},
}

var defaultMenuKeys = map[string]map[int]bool{
	KEYBOARD_G75: {42: true, 98: true, 119: true, 120: true, 121: true},
	KEYBOARD_A75: {42: true, 98: true, 119: true, 120: true, 121: true},
	KEYBOARD_G65: {46: true, 47: true, 48: true, 49: true, 67: true, 68: true, 69: true, 70: true},
	KEYBOARD_G60: {46: true, 47: true, 48: true, 49: true, 67: true, 68: true, 69: true, 70: true},
}

func DefaultFnKeys(model string) map[int]bool {
	if keys, ok := defaultFnKeys[model]; ok {
		return keys
	}
	return nil
}

func DefaultMenuKeys(model string) map[int]bool {
	if keys, ok := defaultMenuKeys[model]; ok {
		return keys
	}
	return nil
}

	var defaultFnActions = map[string]map[int]string{
	KEYBOARD_G60: {
		21: "Tilde\n'", 22: KeyF1, 23: KeyF2, 24: KeyF3, 25: KeyF4,
		26: KeyF5, 27: KeyF6, 28: KeyF7, 29: KeyF8, 30: KeyF9,
		31: KeyF10, 32: KeyF11, 33: KeyF12,
		44: KeyArrUp,
		52: KeyPrint, 53: KeyScrlk, 54: KeyPause,
		64: KeyArrLeft, 65: KeyArrDown, 66: KeyArrRight,
		72: KeyIns, 73: KeyHome, 74: KeyPgUp,
		86: KeyPrev, 87: KeyNext, 88: KeyPlay,
		89: KeyStop, 90: KeyMute,
		91: "\U000F075E", 92: "\U000F075D",
		93: KeyDel, 94: KeyEnd, 95: KeyPgDn,
	},
	KEYBOARD_G65: {
		21: "`~", 22: KeyF1, 23: KeyF2, 24: KeyF3, 25: KeyF4,
		26: KeyF5, 27: KeyF6, 28: KeyF7, 29: KeyF8, 30: KeyF9,
		31: KeyF10, 32: KeyF11, 33: KeyF12,
		35: KeyIns,
		54: KeyPause,
		56: KeyHome,
		77: KeyPrint,
		86: KeyPrev, 87: KeyNext, 88: KeyPlay,
		89: KeyStop, 90: KeyMute,
		91: "\U000F075E", 92: "\U000F075D",
		98: KeyScrlk,
	},
}

var defaultMenuActions = map[string]map[int]string{
	KEYBOARD_G65: {
		46: KeyLightModeSw,
		47: KeyTurbo,
		48: KeyLightColorCycle,
		49: KeyLightModeCycle,
		67: KeyLightLumiDec,
		68: KeyLightLumiInc,
		69: KeyLightSpeedDec,
		70: KeyLightSpeedInc,
	},
	KEYBOARD_G60: {
		46: KeyLightModeSw,
		47: KeyTurbo,
		48: KeyLightColorCycle,
		49: KeyLightModeCycle,
		67: KeyLightLumiDec,
		68: KeyLightLumiInc,
		69: KeyLightSpeedDec,
		70: KeyLightSpeedInc,
	},
	KEYBOARD_G75: {
		42: KeyLightColorCycle,
		98: KeyLightLumiInc,
		119: KeyLightModeNext,
		120: KeyLightLumiDec,
		121: KeyLightModePrev,
	},
	KEYBOARD_A75: {
		42: KeyLightColorCycle,
		98: KeyLightLumiInc,
		119: KeyLightModeNext,
		120: KeyLightLumiDec,
		121: KeyLightModePrev,
	},
}

func DefaultFnActions(model string) map[int]string {
	if actions, ok := defaultFnActions[model]; ok {
		return actions
	}
	return nil
}

func DefaultMenuActions(model string) map[int]string {
	if actions, ok := defaultMenuActions[model]; ok {
		return actions
	}
	return nil
}
