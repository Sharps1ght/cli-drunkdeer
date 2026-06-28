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

func (d *DrunkDeerController) SendCustomColorPacket(colorData []byte, brightness byte, turbo bool) {
	report := BuildCustomColorPacket(colorData, brightness, turbo)
	d.QueuePacket(report)
}

func (d *DrunkDeerController) SendCustomColorData(colors map[int][3]byte, brightness byte, defaultColor [3]byte, turbo bool) {
	const chunkSize = COLORS_PER_PACKET * BYTES_PER_KEY

	model := KEYBOARD_G60
	if d.identity != nil {
		model = d.identity.KeyboardModel
	}

	var flat []byte
	for _, row := range GetLEDGrid(model) {
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
		d.SendCustomColorPacket(flat[offset:end], brightness, turbo)
	}

	// send empty terminator packet (webdriver sends E=6 for G60, last is empty)
	d.SendCustomColorPacket(nil, brightness, turbo)
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

func (d *DrunkDeerController) SendClearRTPData() {
	d.QueuePacket([]byte{0xAA, 0x00, 0x01})
}

func (d *DrunkDeerController) SendRemapData(keys map[int]*RemapKey, layer byte) {
	for chunk := byte(1); chunk <= 14; chunk++ {
		report := BuildRemapChunk(layer, chunk, keys)
		d.QueuePacket(report)
	}
}

func (d *DrunkDeerController) QueuePacket(p []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	packet := packetPool.Get().([]byte)
	n := len(p)
	if n > 63 {
		n = 63
	}
	copy(packet, p)
	for i := n; i < 63; i++ {
		packet[i] = 0
	}

	d.packetWg.Add(1)
	d.packetQueue <- packet
}

func (d *DrunkDeerController) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.packetWg.Wait()
}
