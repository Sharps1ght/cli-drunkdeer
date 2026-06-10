package driver

func (d *DrunkDeerController) SendIdentity() {
	report := BuildIdentity()
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendLEDModeSelect(direction, sequence, speed, brightness, rgb byte) {
	report := BuildLEDModeSelect(direction, sequence, speed, brightness, rgb)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendLEDModeSelectTurbo(direction, sequence, speed, brightness, rgb byte) {
	report := BuildLEDModeSelectTurbo(direction, sequence, speed, brightness, rgb)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendCustomColorPacket(colorData []byte, brightness byte) {
	report := BuildCustomColorPacket(colorData, brightness)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendCustomColorData(colors map[int][3]byte, brightness byte, defaultColor [3]byte) {
	const chunkSize = COLORS_PER_PACKET * BYTES_PER_KEY

	var flat []byte
	for _, row := range G60_LED_GRID {
		for _, idx := range row {
			rgb, ok := colors[idx]
			if !ok {
				rgb = defaultColor
			}
			d.Log("color[%d] = %02x %02x %02x", idx, rgb[0], rgb[1], rgb[2])
			flat = append(flat, CUSTOM_COLOR_PADDING|byte(idx), rgb[0], rgb[1], rgb[2])
		}
	}

	for offset := 0; offset < len(flat); offset += chunkSize {
		end := offset + chunkSize
		if end > len(flat) {
			end = len(flat)
		}
		d.SendCustomColorPacket(flat[offset:end], brightness)
	}

	// send empty terminator packet (webdriver sends E=6 for G60, last is empty)
	d.SendCustomColorPacket(nil, brightness)
}

func (d *DrunkDeerController) SendLEDModeDisable() {
	report := BuildLEDModeDisable()
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendModifyRow(row uint8, keys []byte) {
	report := BuildModifyRowActuation(row, keys)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendRapidTriggerTurbo(rt, turbo bool) {
	report := BuildRapidTriggerTurbo(rt, turbo)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendDownstrokes(row uint8, keys []byte) {
	report := BuildModifyRowDownstroke(row, keys)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendUpstrokes(row uint8, keys []byte) {
	report := BuildModifyRowUpstroke(row, keys)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) QueuePacket(p []byte) {
	var packet []byte
	if len(p) != 63 {
		packet = make([]byte, 63)
		copy(packet[:], p[:])
	} else {
		packet = p
	}

	d.packetWg.Add(1)
	d.packetQueue <- packet
}

func (d *DrunkDeerController) Flush() {
	d.packetWg.Wait()
}
