package driver

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/sstallion/go-hid"
)

func (d *DrunkDeerController) GetIdentity() *DDKeyboardIdentity {
	if d.identity == nil {
		d.SendIdentity()
		for d.identity == nil && !d.shouldClose {
			time.Sleep(1 * time.Millisecond)
		}
	}

	return d.identity
}

func (d *DrunkDeerController) GetActuations() []byte {
	// No way to get actuations from device because it only echoes whatever you throw at it (if it's valid)

	return d.actuations
}

func (d *DrunkDeerController) sendReport(p []byte) {
	d.Log("Sending report: %x", p)
	report := make([]byte, 64)
	report[0] = KEYBOARD_REPORT_ID // Report ID

	if len(p) > 63 {
		p = p[:63]
	}

	copy(report[1:], p)

	_, err := d.device.Write(report)
	if err != nil {
		d.Log("Write error: %v", err)
	}
}

func (d *DrunkDeerController) SetDebug(debug bool) {
	d.debug = debug
}

func (d *DrunkDeerController) Log(str string, v ...interface{}) {
	if str[len(str)-1] != '\n' {
		str += "\n"
	}

	if d.debug {
		fmt.Printf(color.HiGreenString("[DEBUG] ")+str, v...)
	}
}

func (d *DrunkDeerController) LoadActuations(actuations []byte) {
	if len(actuations) != LAYOUT_SIZE {
		panic("Actuations length does not match keyboard layout length")
	}

	d.actuations = actuations

	for i := 0; i < len(d.actuations); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.actuations) {
			end = len(d.actuations)
		}
		row := d.actuations[i:end]

		d.SendModifyRow(uint8(i/KEYS_PER_ROW), row)
	}
}

func (d *DrunkDeerController) LoadDownstrokes(downstrokes []byte) {
	if len(downstrokes) != LAYOUT_SIZE {
		panic("Downstrokes length does not match keyboard layout length")
	}

	d.downstrokes = downstrokes

	for i := 0; i < len(d.downstrokes); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.downstrokes) {
			end = len(d.downstrokes)
		}

		row := d.downstrokes[i:end]
		rowIndex := uint8(i / KEYS_PER_ROW)
		d.SendDownstrokes(rowIndex, row)
	}
}

func (d *DrunkDeerController) LoadUpstrokes(upstrokes []byte) {
	if len(upstrokes) != LAYOUT_SIZE {
		panic("Upstrokes length does not match keyboard layout length")
	}

	d.upstrokes = upstrokes

	for i := 0; i < len(d.upstrokes); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.upstrokes) {
			end = len(d.upstrokes)
		}

		row := d.upstrokes[i:end]
		rowIndex := uint8(i / KEYS_PER_ROW)
		d.SendUpstrokes(rowIndex, row)
	}
}

// #region Modifiers
func (d *DrunkDeerController) ModifyActuationsByNames(names []string, actuations byte) {
	for _, name := range names {
		index := GetIndexByKey(name, KEYBOARD_A75)
		if index != -1 {
			d.actuations[index] = actuations
		}
	}
}

func (d *DrunkDeerController) ModifyActuationsByIndexes(indexes []int, actuation byte) {
	for _, index := range indexes {
		if index >= 0 && index < len(d.actuations) {
			d.actuations[index] = actuation
		}
	}
}

func (d *DrunkDeerController) ModifyAllActuations(actuation byte) {
	for i := range d.actuations {
		d.actuations[i] = actuation
	}
}

// #endregion

func (d *DrunkDeerController) WriteDefaults() {
	actuations := make([]byte, KEYS_PER_ROW)
	for i := 0; i < KEYS_PER_ROW; i++ {
		actuations[i] = DEFAULT_ACTUATION
	}

	d.Log("Writing defaults")
	d.SendLEDModeDisable()
	d.SendRapidTriggerTurbo(false, false)

	for i := 0; i < len(d.actuations); i += KEYS_PER_ROW {
		end := i + KEYS_PER_ROW
		if end > len(d.actuations) {
			end = len(d.actuations)
		}
		rowIndex := uint8(i / KEYS_PER_ROW)

		row := d.actuations[i:end]
		d.SendModifyRow(rowIndex, row)

		downstrokes := d.downstrokes[i:end]
		d.SendDownstrokes(rowIndex, downstrokes)

		upstrokes := d.upstrokes[i:end]
		d.SendUpstrokes(rowIndex, upstrokes)
	}

	d.Log("Defaults written")
}

func (d *DrunkDeerController) Close() error {
	var closeErr error
	d.closeOnce.Do(func() {
		d.shouldClose = true

		close(d.packetQueue)
		close(d.packetChan)

		done := make(chan struct{})
		go func() {
			d.wg.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			closeErr = fmt.Errorf("timeout waiting for goroutines to exit")
		}
	})
	return closeErr
}

func NewDrunkDeerController(device *hid.Device) *DrunkDeerController {
	controller := &DrunkDeerController{
		device:      device,
		packetChan:  make(chan DDPacket),
		packetQueue: make(chan []byte, 10),
		Light:       &DDLight{},
	}

	controller.actuations = make([]byte, LAYOUT_SIZE)
	controller.downstrokes = make([]byte, LAYOUT_SIZE)
	controller.upstrokes = make([]byte, LAYOUT_SIZE)

	for i := 0; i < LAYOUT_SIZE; i++ {
		controller.actuations[i] = DEFAULT_ACTUATION
		controller.downstrokes[i] = 0x00
		controller.upstrokes[i] = 0x00
	}

	go controller.drunkDeerMessageReceiver()
	controller.wg.Add(1)

	// Start a goroutine to read from the device and send packets to the channel
	go func() {
		defer controller.wg.Done()
		defer func() { recover() }()
		for {
			if controller.shouldClose {
				return // Exit when shouldClose is set
			}

			buf := make([]byte, 64)
			n, err := device.Read(buf)
			if err != nil {
				return // Exit if reading fails
			}

			if n > 0 {
				if buf[0] != KEYBOARD_REPORT_ID {
					if controller.debug {
						fmt.Printf(color.HiYellowString("[DEBUG ALT REPORT ID 0x%02x] %x\n"), buf[0], buf[1:n])
					}
					continue
				}

				packet := DDPacket{
					Packet: buf[1],
					Data:   buf[2:n],
				}
				select {
				case controller.packetChan <- packet:
				case <-time.After(100 * time.Millisecond):
				}
			}
		}
	}()
	controller.wg.Add(1)

	go controller.drunkDeerReporter()
	controller.wg.Add(1)

	return controller
}

func (d *DrunkDeerController) drunkDeerReporter() {
	defer d.wg.Done()
	for {
		select {
		case p, ok := <-d.packetQueue:
			if !ok {
				return
			}
			d.sendReport(p)
			d.packetWg.Done()
			time.Sleep(5 * time.Millisecond)
		case <-time.After(100 * time.Millisecond):
			if d.shouldClose {
				return
			}
		}
	}
}

func (d *DrunkDeerController) drunkDeerMessageReceiver() {
	i := 0
	defer d.wg.Done()
	for p := range d.packetChan {
		data := bytes.NewBuffer(p.Data)
		expectedValue := data.Next(1)
		i += 1

		d.Log("%d Packet received: %x", i, p.Data)
		switch p.Packet {
		case PACKET_IDENTITY:
			// fmt.Printf("Identity packet received: %x\n", p.Data)
			theNullByte, _ := data.ReadByte()
			data.ReadByte()
			modelBytes := data.Next(3)

			if theNullByte != 0x00 {
				break
			}

			if expectedValue[0] != 0x02 {
				d.Log("Unknown expected value: %x", expectedValue[0])
			}

			versionBytes := data.Next(2)
			model, keyboardType := DetectKeyboardModel(modelBytes)
			version := int(uint16(versionBytes[0]) | uint16(versionBytes[1])<<8)
			ident := DDKeyboardIdentity{
				KeyboardModel:   model,
				KeyboardType:    uint8(keyboardType),
				FirmwareVersion: fmt.Sprintf("0.0%v", version),
				RapidTrigger:    p.Data[15] != 0,
				Turbo:           p.Data[14] != 0,
			}
			d.identity = &ident
		case PACKET_LEDMODESEL:
			d.Light.Direction = p.Data[2]
			d.Light.Sequence = p.Data[3]
			d.Light.Speed = p.Data[4]
			d.Light.Brightness = p.Data[5]
		case PACKET_TURBORT:
			d.Log("Turbo packet received: %x", p.Data)
			d.turbo = p.Data[6] != 0
			d.rapidTrigger = p.Data[7] != 0
		case PACKET_MODIFYKEY:
			d.Log("Modify key packet received: %x", p.Data)
		case PACKET_KEYTRACKING:
		default:
			d.Log("Unknown packet type: %x", p.Packet)
			d.Log("Data: %x", p.Data)
		}
	}
}
