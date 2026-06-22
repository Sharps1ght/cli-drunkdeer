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
		21: "`~", 22: "F1", 23: "F2", 24: "F3", 25: "F4",
		26: "F5", 27: "F6", 28: "F7", 29: "F8", 30: "F9",
		31: "F10", 32: "F11", 33: "F12",
		44: "ARR_UP",
		52: "PrtSc", 53: "scroll", 54: "Pause",
		64: "ARR_LEFT", 65: "ARR_DOWN", 66: "ARR_RIGHT",
		72: "Ins", 73: "Home", 74: "PgUp",
		86: "PREV", 87: "NEXT", 88: "PLAY",
		89: "STOP", 90: "MUTE",
		91: "\U000F075E", 92: "\U000F075D",
		93: "Del", 94: "End", 95: "PgDn",
	},
	KEYBOARD_G65: {
		21: "`~", 22: "F1", 23: "F2", 24: "F3", 25: "F4",
		26: "F5", 27: "F6", 28: "F7", 29: "F8", 30: "F9",
		31: "F10", 32: "F11", 33: "F12",
		35: "Ins",
		54: "Pause",
		56: "Home",
		77: "PrtSc",
		86: "PREV", 87: "NEXT", 88: "PLAY",
		89: "STOP", 90: "MUTE",
		91: "\U000F075E", 92: "\U000F075D",
		98: "scroll",
	},
}

var defaultMenuActions = map[string]map[int]string{
	KEYBOARD_G65: {
		46: "Light",
		47: "Turbo",
	},
	KEYBOARD_G60: {
		46: "Light",
		47: "Turbo",
	},
	KEYBOARD_G75: {
		42: "Color",
		98: "Lumi+",
		119: "Mode\u2190",
		120: "Lumi\u2212",
		121: "Mode\u2192",
	},
	KEYBOARD_A75: {
		42: "Color",
		98: "Lumi+",
		119: "Mode\u2190",
		120: "Lumi\u2212",
		121: "Mode\u2192",
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
