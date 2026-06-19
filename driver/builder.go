package driver

func BuildIdentity() []byte {
	return []byte{
		PACKET_IDENTITY,
		0x02,
	}
}

func BuildLEDModeSelect(direction, sequence, speed, brightness, rgb byte) []byte {
	return []byte{
		PACKET_LEDMODESEL,
		0x01,
		0x00,
		direction,
		sequence,
		speed,
		brightness,
		rgb,
	}
}

func BuildLEDModeDisable() []byte {
	return []byte{
		PACKET_LEDMODESEL,
		0x01,
		0x00,
		0x00,
		SEQUENCE_OFF,
		0x00,
		0x00,
		0x00,
	}
}

func BuildLEDModeSelectTurbo(direction, sequence, speed, brightness, rgb byte) []byte {
	report := BuildLEDModeSelect(direction, sequence, speed, brightness, rgb)
	report[2] = 0x01

	return report
}

func BuildRapidTriggerTurbo(rt, turbo bool) []byte {
	return []byte{
		PACKET_TURBORT,
		0x00,
		0x1E,
		0x01,
		0x00,
		0x00,
		0x01,
		BoolToByte(turbo),
		BoolToByte(rt),
	}
}

func BuildKeyTracking(track bool) []byte {
	return []byte{
		PACKET_MODIFYKEY,
		0x03,
		BoolToByte(track),
	}
}

func BuildCustomColorPacket(colorData []byte, brightness byte, turbo bool) []byte {
	report := make([]byte, 63)
	report[0] = PACKET_LEDMODESEL
	report[1] = 0x01
	if turbo {
		report[2] = 0x01 // turbo flag (1 = turbo custom mode)
	} else {
		report[2] = 0x00 // turbo flag (0 = non-turbo custom mode)
	}
	report[3] = 0x00
	report[4] = SEQUENCE_CUSTOM
	report[5] = 0x06 // data type: custom color
	report[6] = brightness
	report[7] = 0xFF // header terminator
	copy(report[8:], colorData)
	report[8+len(colorData)] = 0xFF

	return report
}

// Row 0 is the first row, row 1 is the second row, and row 2 is the third row
func BuildModifyRow(row uint8, keys []byte, defaultValue byte) []byte {
	report := make([]byte, 63)
	report[0] = PACKET_MODIFYKEY
	report[1] = 0x01
	report[2] = 0x00
	report[3] = row

	maxKeys := 59
	if row == 2 {
		maxKeys = 8
	}
	if len(keys) > maxKeys {
		keys = keys[:maxKeys]
	}

	// Fill report starting at offset 4
	for i := 0; i < maxKeys; i++ {
		if i < len(keys) {
			report[i+4] = keys[i]
		} else {
			report[i+4] = defaultValue
		}
	}

	return report
}

func BuildModifyRowActuation(row uint8, keys []byte) []byte {
	report := BuildModifyRow(row, keys, DEFAULT_ACTUATION)
	report[1] = 0x01

	return report
}

func BuildModifyRowDownstroke(row uint8, keys []byte) []byte {
	report := BuildModifyRow(row, keys, 0x00) // 0x00 is 0.0mm
	report[1] = 0x04

	return report
}

func BuildModifyRowUpstroke(row uint8, keys []byte) []byte {
	report := BuildModifyRow(row, keys, 0x00) // 0x00 is 0.0mm
	report[1] = 0x05

	return report
}
