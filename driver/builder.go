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

func BuildCustomColorPacket(colorData []byte, brightness byte) []byte {
	report := make([]byte, 63)
	report[0] = PACKET_LEDMODESEL
	report[1] = 0x01
	report[2] = 0x00 // turbo flag (0 = non-turbo custom mode)
	report[3] = 0x00
	report[4] = SEQUENCE_CUSTOM
	report[5] = 0x06 // data type: custom color
	report[6] = brightness
	report[7] = 0xFF // header terminator
	copy(report[8:], colorData)
	// terminator right after color data, matching webdriver's sendTurboLedModeData
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

	if len(keys) > 59 {
		keys = keys[:59]
	}

	if len(keys) > 8 && row == 2 {
		keys = keys[:8]
	}

	if len(keys) < 59 && row != 2 {
		for i := len(keys); i < 59; i++ {
			keys = append(keys, defaultValue)
		}
	}

	// I'm paranoid about drunkdeer firmware, so let's make sure we have the right amount of keys
	if len(keys) < 8 && row == 2 {
		for i := len(keys); i < 8; i++ {
			keys = append(keys, defaultValue)
		}
	}

	for i := 0; i < len(keys); i++ {
		report[i+4] = keys[i]
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
