package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/2xxn/cli-drunkdeer/driver"
	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
	"github.com/sstallion/go-hid"
)

type App struct {
	keyboardIndex int
	device        *hid.Device
	ledDevice     *hid.Device
	controller    *driver.DrunkDeerController
	profilePath   string
	args          Args
}

func (a *App) parseArgs() {
	arg.MustParse(&a.args)
	debug = a.args.Debug
}

func (a *App) setupProfilePath() {
	if a.profilePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		handleError("Error getting home directory", err)

		a.profilePath = homeDir + a.profilePath[1:]
	}

	if err := os.MkdirAll(a.profilePath, 0755); err != nil {
		handleError("Error creating profile directory", err)
	}
}

func (a *App) setupDevice() {
	var err error

	devices := FindDrunkDeerDevices()
	if len(devices) == 0 {
		handleError("Error", fmt.Errorf("no devices found"))
	}
	if a.keyboardIndex >= len(devices) {
		a.keyboardIndex = 0
	}

	a.device, err = hid.OpenPath(devices[a.keyboardIndex].Path)
	handleError("Error opening main device:", err)
	DEBUG("Main device opened: %s", devices[a.keyboardIndex].Path)

	if len(devices) > a.keyboardIndex+1 {
		a.ledDevice, err = hid.OpenPath(devices[a.keyboardIndex+1].Path)
		if err == nil {
			DEBUG("LED device opened: %s", devices[a.keyboardIndex+1].Path)
		}
	}
	if a.ledDevice == nil && a.keyboardIndex > 0 {
		a.ledDevice, err = hid.OpenPath(devices[0].Path)
		if err == nil {
			DEBUG("LED device opened (fallback): %s", devices[0].Path)
		}
	}

	a.controller = driver.NewDrunkDeerController(a.device)
	a.controller.GetIdentity()
	DEBUG("Created controller")

	if debug {
		a.controller.SetDebug(true)
	}
}

func (a *App) sendLEDReport(p []byte) {
	dev := a.ledDevice
	if dev == nil {
		dev = a.device
	}
	report := make([]byte, 64)
	report[0] = driver.KEYBOARD_REPORT_ID
	if len(p) > 63 {
		p = p[:63]
	}
	copy(report[1:], p)
	DEBUG("Sending LED report: %x", report)
	_, err := dev.Write(report)
	if err != nil {
		DEBUG("LED write error: %v", err)
	}
}

func (a *App) cleanup() {
	if a.controller != nil {
		a.controller.Close()
	}
	if a.ledDevice != nil {
		a.ledDevice.Close()
	}
	if a.device != nil {
		a.device.Close()
	}
}

func (a *App) run() {
	switch {
	case a.args.Reset:
		a.handleReset()
	case a.args.Load != "":
		a.handleLoadProfile()
	default:
		a.showHelp()
	}
}

func (a *App) handleReset() {
	color.HiRed("Resetting device to default settings")
	a.controller.WriteDefaults()
	a.controller.SendLEDModeDisable()
	a.controller.Flush()
	color.White("Reset complete")
}

func (a *App) handleLoadProfile() {
	config := a.getConfig(a.args.Load)
	if config.Model != "" && config.Model != a.controller.GetIdentity().KeyboardModel {
		color.HiRed("Profile model does not match device model (expected %s, got %s)",
			a.controller.GetIdentity().KeyboardModel, config.Model)
		os.Exit(1)
	}

	DEBUG("Model: %v | Turbo: %v | RT: %v | Default actuation: %v",
		config.Model, config.Turbo, config.RapidTrigger.Enabled, config.DefaultActuation)

	actuations, downstrokes, upstrokes := a.prepareKeySettings(config)

	if !config.Light.Enabled {
		config.Light.Sequence = driver.SEQUENCE_OFF
	}

	a.configureLights(config)
	a.applySettings(config, actuations, downstrokes, upstrokes)
	a.applyRemap(config)

	color.White("Loaded %s%s%s",
		color.GreenString(a.args.Load),
		color.WhiteString(" for "),
		color.HiBlueString("DrunkDeer %s", config.Model))
	DEBUG("Profile loaded")
}

func (a *App) prepareKeySettings(config *Config) ([]byte, []byte, []byte) {
	actuations := make([]byte, len(driver.KEYBOARD_LAYOUT))
	downstrokes := make([]byte, len(driver.KEYBOARD_LAYOUT))
	upstrokes := make([]byte, len(driver.KEYBOARD_LAYOUT))

	defaultAct := driver.ActuationFloatToByte(config.DefaultActuation)
	defaultDS := driver.ActuationFloatToByte(config.RapidTrigger.DefaultDownstroke)
	defaultUS := driver.ActuationFloatToByte(config.RapidTrigger.DefaultUpstroke)

	for i := range actuations {
		actuations[i] = defaultAct
		downstrokes[i] = defaultDS
		upstrokes[i] = defaultUS
	}

	for key, value := range config.ActuationPoints {
		i := driver.GetIndexByKey(key)
		actuations[i] = driver.ActuationFloatToByte(value)
	}

	for key, value := range config.RapidTriggers {
		i := driver.GetIndexByKey(key)
		downstrokes[i] = driver.ActuationFloatToByte(value[0])
		upstrokes[i] = driver.ActuationFloatToByte(value[1])
	}

	return actuations, downstrokes, upstrokes
}

func parseHexColor(s string) [3]byte {
	var rgb [3]byte
	s = strings.TrimPrefix(s, "#")
	if len(s) == 6 {
		b, err := hex.DecodeString(s)
		if err == nil && len(b) == 3 {
			rgb = [3]byte{b[0], b[1], b[2]}
		}
	}
	return rgb
}

func (a *App) configureLights(config *Config) {
	colors := make(map[int][3]byte)
	for keyName, hexStr := range config.Light.Colors {
		idx := driver.GetIndexByKey(keyName)
		if config.Model == "G60" {
			if g60idx := driver.GetG60IndexByKey(keyName); g60idx != -1 {
				idx = g60idx
			}
		}
		if idx != -1 {
			colors[idx] = parseHexColor(hexStr)
		}
	}

	defaultCol := [3]byte{0, 0, 0}
	var colIdx byte
	if config.Light.Color != nil {
		if config.Light.Color.Fill != "" {
			defaultCol = parseHexColor(config.Light.Color.Fill)
		}
		colIdx = byte(config.Light.Color.Index)
	}

	a.controller.Light = &driver.DDLight{
		Sequence:     byte(config.Light.Sequence),
		Speed:        byte(config.Light.Speed),
		Direction:    byte(config.Light.Direction),
		Brightness:   byte(config.Light.Brightness),
		Color:        colIdx,
		Colors:       colors,
		DefaultColor: defaultCol,
	}
}

func (a *App) applySettings(config *Config, actuations, downstrokes, upstrokes []byte) {
	a.controller.SendRapidTriggerTurbo(config.RapidTrigger.Enabled, config.Turbo)
	a.controller.LoadActuations(actuations)
	a.controller.LoadDownstrokes(downstrokes)
	a.controller.LoadUpstrokes(upstrokes)
	a.controller.Flush()

	seq := a.controller.Light.Sequence
	if seq == 0 || seq > driver.SEQUENCE_CUSTOM {
		seq = driver.SEQUENCE_WAVE
	}
	clr := a.controller.Light.Color
	if clr == 0 {
		clr = driver.COLOR_RED
	}

	if seq == driver.SEQUENCE_CUSTOM {
		DEBUG("Sending custom color data (%d keys, default fill %02x%02x%02x)",
			len(a.controller.Light.Colors),
			a.controller.Light.DefaultColor[0], a.controller.Light.DefaultColor[1], a.controller.Light.DefaultColor[2])
		a.controller.SendCustomColorData(a.controller.Light.Colors, a.controller.Light.Brightness, a.controller.Light.DefaultColor)
	} else {
		a.controller.SendLEDModeSelect(
			0x00,
			seq,
			a.controller.Light.Speed,
			a.controller.Light.Brightness,
			clr,
		)
	}
	a.controller.Flush()
}

func (a *App) applyRemap(config *Config) {
	model := a.controller.GetIdentity().KeyboardModel
	layers := []struct {
		layer byte
		entries map[string]string
	}{
		{1, config.Remap.Default},
		{2, config.Remap.Fn1},
		{3, config.Remap.Fn2},
	}

	for _, l := range layers {
		if len(l.entries) == 0 {
			continue
		}
		keys := make(map[int]*driver.RemapKey)
		for physKey, action := range l.entries {
			idx := driver.GetRemapIndexByKey(physKey, model)
			if idx == -1 {
				DEBUG("Remap: unknown physical key %q (layer %d)", physKey, l.layer)
				continue
			}
			rk, ok := driver.GetRemapAction(action)
			if !ok {
				DEBUG("Remap: unknown action %q", action)
				continue
			}
			keys[idx] = &rk
			DEBUG("Remap: %s(%d) -> %s (cmd=%02x code=%02x type=%d)", physKey, idx, action, rk.KeyCmd, rk.KeyCode, rk.KeyType)
		}
		if len(keys) > 0 {
			a.controller.SendRemapData(keys, l.layer)
		}
	}
	a.controller.Flush()
}

func (a *App) showHelp() {
	if a.args.Command != "" {
		color.HiRed("Unknown command: %s\n", color.HiWhiteString(a.args.Command))
	}

	color.HiBlue("List of commands")
	color.White("For descriptions, run: drunkdeer --help")
	color.HiWhite("  - drunkdeer import <url/path>")
	color.HiWhite("  - drunkdeer load <profile>")
	color.HiWhite("  - drunkdeer save <profile>")
	color.HiWhite("  - drunkdeer profiles")
	color.HiWhite("  - drunkdeer reset")
	color.HiWhite("  - drunkdeer list")
	color.HiWhite("  - drunkdeer version")
}
