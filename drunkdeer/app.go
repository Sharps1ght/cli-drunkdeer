package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	if a.device != nil {
		a.device.Close()
	}
	if a.ledDevice != nil {
		a.ledDevice.Close()
	}
	if a.controller != nil {
		a.controller.Close()
	}
}

func (a *App) run() {
	switch {
	case a.args.Reset:
		a.handleReset()
	case a.args.Load != "":
		a.handleLoadProfile()
	case a.args.Command == "set":
		a.handleSet()
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
	t0 := time.Now()
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
	DEBUG("prepareKeySettings + configureLights: %v", time.Since(t0))
	t1 := time.Now()

	a.applySettings(config, actuations, downstrokes, upstrokes)
	DEBUG("applySettings: %v", time.Since(t1))
	t2 := time.Now()

	a.applyRemap(config)
	DEBUG("applyRemap: %v", time.Since(t2))

	color.White("Loaded %s%s%s",
		color.GreenString(a.args.Load),
		color.WhiteString(" for "),
		color.HiBlueString("DrunkDeer %s", config.Model))
	DEBUG("Total load time: %v", time.Since(t0))
}

func (a *App) prepareKeySettings(config *Config) ([]byte, []byte, []byte) {
	actuations := make([]byte, driver.LAYOUT_SIZE)
	downstrokes := make([]byte, driver.LAYOUT_SIZE)
	upstrokes := make([]byte, driver.LAYOUT_SIZE)

	defaultAct := driver.ActuationFloatToByte(config.DefaultActuation)
	defaultDS := driver.ActuationFloatToByte(config.RapidTrigger.DefaultDownstroke)
	defaultUS := driver.ActuationFloatToByte(config.RapidTrigger.DefaultUpstroke)

	for i := range actuations {
		actuations[i] = defaultAct
		downstrokes[i] = defaultDS
		upstrokes[i] = defaultUS
	}

	for key, value := range config.ActuationPoints {
		i := driver.GetIndexByKey(key, config.Model)
		actuations[i] = driver.ActuationFloatToByte(value)
	}

	for key, value := range config.RapidTriggers {
		i := driver.GetIndexByKey(key, config.Model)
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
		idx := driver.GetIndexByKey(keyName, config.Model)
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

	turboDefaultCol := defaultCol
	if config.Light.TurboColor != nil {
		if config.Light.TurboColor.Fill != "" {
			turboDefaultCol = parseHexColor(config.Light.TurboColor.Fill)
		}
	}

	a.controller.Light = &driver.DDLight{
		Sequence:     byte(config.Light.Sequence),
		Speed:        byte(config.Light.Speed),
		Direction:    byte(config.Light.Direction),
		Brightness:   byte(config.Light.Brightness),
		Color:        colIdx,
		Colors:       colors,
		DefaultColor: defaultCol,

		TurboDefaultColor: turboDefaultCol,
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
		light := a.controller.Light

		// webdriver sends mode select (byte[7]=0xFF = "data follows") before color chunks
		a.controller.SendLEDModeSelect(0x00, driver.SEQUENCE_CUSTOM, light.Speed, light.Brightness, 0xFF)
		a.controller.Flush()
		time.Sleep(10 * time.Millisecond)

		if config.Turbo {
			DEBUG("Sending turbo custom color data (uniform, fill %02x%02x%02x)",
				light.TurboDefaultColor[0], light.TurboDefaultColor[1], light.TurboDefaultColor[2])
			a.controller.SendCustomColorData(
				make(map[int][3]byte),
				light.Brightness,
				light.TurboDefaultColor,
				true,
			)
		} else {
			DEBUG("Sending custom color data (%d keys, default fill %02x%02x%02x)",
				len(light.Colors),
				light.DefaultColor[0], light.DefaultColor[1], light.DefaultColor[2])
			a.controller.SendCustomColorData(light.Colors, light.Brightness, light.DefaultColor, false)
		}
		a.controller.Flush()
		return
	} else {
		// send mode select (own flush + delay so actuation packets don't reset it)
		a.controller.SendLEDModeSelect(
			0x00,
			seq,
			a.controller.Light.Speed,
			a.controller.Light.Brightness,
			clr,
		)
		a.controller.Flush()
		time.Sleep(10 * time.Millisecond)
	}
}

func (a *App) applyRemap(config *Config) {
	model := a.controller.GetIdentity().KeyboardModel
	layers := []struct {
		layer byte
		entries map[string]string
	}{
		{1, config.Remap.Default},
		{2, config.Remap.Fn},
		{3, config.Remap.Menu},
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

func (a *App) loadSetState() *Config {
	path := filepath.Join(a.profilePath, ".state.json")
	data, err := os.ReadFile(path)
	if err != nil {
		model := a.controller.GetIdentity().KeyboardModel
		return &Config{
			Model:            model,
			DefaultActuation: 2.0,
			RapidTrigger: RapidTriggerSettings{
				DefaultDownstroke: 2.0,
				DefaultUpstroke:   2.0,
			},
			Light: LightSettings{
				Enabled:    true,
				Sequence:   driver.SEQUENCE_CUSTOM,
				Brightness: 9,
				Speed:      9,
				Colors:     map[string]string{},
			},
			ActuationPoints: map[string]float32{},
			RapidTriggers:   map[string][2]float32{},
			Remap:           RemapSettings{},
		}
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	if cfg.Light.Colors == nil {
		cfg.Light.Colors = map[string]string{}
	}
	if cfg.ActuationPoints == nil {
		cfg.ActuationPoints = map[string]float32{}
	}
	if cfg.RapidTriggers == nil {
		cfg.RapidTriggers = map[string][2]float32{}
	}
	return &cfg
}

func (a *App) saveSetState(cfg *Config) {
	path := filepath.Join(a.profilePath, ".state.json")
	data, _ := json.MarshalIndent(cfg, "", "    ")
	os.WriteFile(path, data, 0644)
}

func (a *App) handleSet() {
	args := a.args.CmdValue
	if len(args) < 2 {
		color.HiRed("Usage: drunkdeer set <section> <key> [value...]")
		os.Exit(1)
	}

	section := args[0]
	key := args[1]
	values := args[2:]
	model := a.controller.GetIdentity().KeyboardModel
	state := a.loadSetState()

	switch section {
	case "turbo":
		v := strings.ToLower(key) == "true" || key == "1"
		state.Turbo = v
		a.controller.SendRapidTriggerTurbo(a.controller.GetIdentity().RapidTrigger, v)
		a.controller.Flush()
		color.HiGreen("Turbo = %v", v)

	case "defaultActuation":
		v, err := strconv.ParseFloat(key, 32)
		if err != nil {
			color.HiRed("Invalid actuation value: %s", key)
			os.Exit(1)
		}
		state.DefaultActuation = float32(v)
		b := driver.ActuationFloatToByte(float32(v))
		for row := uint8(0); row < 3; row++ {
			keys := make([]byte, 59)
			for i := range keys {
				keys[i] = b
			}
			a.controller.SendModifyRow(row, keys)
		}
		a.controller.Flush()
		color.HiGreen("Default actuation = %.1fmm", v)

	case "light":
		switch key {
		case "color":
			if len(values) < 1 {
				color.HiRed("Usage: drunkdeer set light color <hex>")
				os.Exit(1)
			}
			state.Light.Color = &ColorSetting{Fill: values[0]}
			sendCustomColor(state.Light.Colors, values[0], a.controller)
			color.HiGreen("Default light color = %s", values[0])

		case "turboColor":
			if len(values) < 1 {
				color.HiRed("Usage: drunkdeer set light turboColor <hex>")
				os.Exit(1)
			}
			rgb := parseHexColor(values[0])
			a.controller.Light.TurboDefaultColor = rgb
			state.Light.TurboColor = &ColorSetting{Fill: values[0]}
			br := setDefaultBrightness(a.controller)
			a.controller.SendCustomColorData(
				make(map[int][3]byte), br, rgb, true,
			)
			a.controller.Flush()
			color.HiGreen("Turbo light color = %s", values[0])

		case "colors":
			if len(values) < 2 {
				color.HiRed("Usage: drunkdeer set light colors <key> <hex>")
				os.Exit(1)
			}
			for _, name := range splitKeys(values[0]) {
				state.Light.Colors[strings.TrimSpace(name)] = values[1]
			}
			// use existing default: controller > state > per-key color
			defaultFill := ""
			if a.controller.Light.DefaultColor != [3]byte{} {
				c := a.controller.Light.DefaultColor
				defaultFill = fmt.Sprintf("#%02x%02x%02x", c[0], c[1], c[2])
			} else if state.Light.Color != nil && state.Light.Color.Fill != "" {
				defaultFill = state.Light.Color.Fill
			}
			if defaultFill == "" {
				defaultFill = values[1]
			}
			sendCustomColor(state.Light.Colors, defaultFill, a.controller)
			color.HiGreen("Light colors %s = %s", values[0], values[1])

		case "brightness":
			if len(values) < 1 {
				color.HiRed("Usage: drunkdeer set light brightness <0-9>")
				os.Exit(1)
			}
			v, err := strconv.Atoi(values[0])
			if err != nil || v < 0 || v > 9 {
				color.HiRed("Invalid brightness (0-9): %s", values[0])
				os.Exit(1)
			}
			state.Light.Brightness = v
			color.HiGreen("Brightness = %d", v)

		case "sequence":
			if len(values) < 1 {
				color.HiRed("Usage: drunkdeer set light sequence <0-19>")
				os.Exit(1)
			}
			v, err := strconv.Atoi(values[0])
			if err != nil || v < 0 || v > 19 {
				color.HiRed("Invalid sequence (0-19): %s", values[0])
				os.Exit(1)
			}
			state.Light.Sequence = v
			br := byte(state.Light.Brightness)
			if br == 0 {
				br = 9
			}
			a.controller.SendLEDModeSelect(0x00, byte(v), byte(state.Light.Speed), br, 0x00)
			a.controller.Flush()
			color.HiGreen("Light sequence = %d", v)

		case "speed":
			if len(values) < 1 {
				color.HiRed("Usage: drunkdeer set light speed <0-9>")
				os.Exit(1)
			}
			v, err := strconv.Atoi(values[0])
			if err != nil || v < 0 || v > 9 {
				color.HiRed("Invalid speed (0-9): %s", values[0])
				os.Exit(1)
			}
			state.Light.Speed = v
			br := byte(state.Light.Brightness)
			if br == 0 {
				br = 9
			}
			a.controller.SendLEDModeSelect(0x00, byte(state.Light.Sequence), byte(v), br, 0x00)
			a.controller.Flush()
			color.HiGreen("Light speed = %d", v)

		default:
			color.HiRed("Unknown light field: %s", key)
			color.White("Valid fields: color, turboColor, colors, brightness, sequence, speed")
			os.Exit(1)
		}

	case "actuationPoints":
		if len(values) < 1 {
			color.HiRed("Usage: drunkdeer set actuationPoints <key> <value>")
			os.Exit(1)
		}
		v, err := strconv.ParseFloat(values[0], 32)
		if err != nil {
			color.HiRed("Invalid actuation value: %s", values[0])
			os.Exit(1)
		}
		for _, name := range splitKeys(key) {
			state.ActuationPoints[strings.TrimSpace(name)] = float32(v)
		}
		b := driver.ActuationFloatToByte(float32(v))
		rows := [3][]byte{{}, {}, {}}
		for _, name := range splitKeys(key) {
			name = strings.TrimSpace(name)
			idx, ok := driver.ResolveKeyToIndex(name, model)
			if !ok {
				color.HiRed("Unknown key: %s", name)
				os.Exit(1)
			}
			row := idx / 59
			col := idx % 59
			if rows[row] == nil {
				rows[row] = make([]byte, 59)
			}
			rows[row][col] = b
		}
		for r := uint8(0); r < 3; r++ {
			if rows[r] != nil {
				a.controller.SendModifyRow(r, rows[r])
			}
		}
		a.controller.Flush()
		color.HiGreen("Actuation %s = %.1fmm", key, v)

	case "rapidTrigger":
		if key == "enabled" {
			v := len(values) > 0 && (strings.ToLower(values[0]) == "true" || values[0] == "1")
			state.RapidTrigger.Enabled = v
			a.controller.SendRapidTriggerTurbo(v, a.controller.GetIdentity().Turbo)
			a.controller.Flush()
			color.HiGreen("Rapid Trigger = %v", v)
		} else if key == "defaultDownstroke" {
			if len(values) < 1 {
				color.HiRed("Usage: drunkdeer set rapidTrigger defaultDownstroke <value>")
				os.Exit(1)
			}
			v, _ := strconv.ParseFloat(values[0], 32)
			state.RapidTrigger.DefaultDownstroke = float32(v)
			color.HiGreen("Rapid trigger default downstroke = %.1fmm", v)
		} else if key == "defaultUpstroke" {
			if len(values) < 1 {
				color.HiRed("Usage: drunkdeer set rapidTrigger defaultUpstroke <value>")
				os.Exit(1)
			}
			v, _ := strconv.ParseFloat(values[0], 32)
			state.RapidTrigger.DefaultUpstroke = float32(v)
			color.HiGreen("Rapid trigger default upstroke = %.1fmm", v)
		} else {
			if len(values) < 2 {
				color.HiRed("Usage: drunkdeer set rapidTrigger <key> <downstroke> <upstroke>")
				os.Exit(1)
			}
			down, err1 := strconv.ParseFloat(values[0], 32)
			up, err2 := strconv.ParseFloat(values[1], 32)
			if err1 != nil || err2 != nil {
				color.HiRed("Invalid rapid trigger values: %s %s", values[0], values[1])
				os.Exit(1)
			}
			for _, name := range splitKeys(key) {
				state.RapidTriggers[strings.TrimSpace(name)] = [2]float32{float32(down), float32(up)}
			}
			ds := driver.ActuationFloatToByte(float32(down))
			us := driver.ActuationFloatToByte(float32(up))
			rows := [3][2][]byte{{}, {}, {}}
			for _, name := range splitKeys(key) {
				name = strings.TrimSpace(name)
				idx, ok := driver.ResolveKeyToIndex(name, model)
				if !ok {
					color.HiRed("Unknown key: %s", name)
					os.Exit(1)
				}
				row := idx / 59
				col := idx % 59
				if rows[row][0] == nil {
					rows[row][0] = make([]byte, 59)
					rows[row][1] = make([]byte, 59)
				}
				rows[row][0][col] = ds
				rows[row][1][col] = us
			}
			for r := uint8(0); r < 3; r++ {
				if rows[r][0] != nil {
					a.controller.SendDownstrokes(r, rows[r][0])
					a.controller.SendUpstrokes(r, rows[r][1])
				}
			}
			a.controller.Flush()
			color.HiGreen("Rapid Trigger %s = %.1fmm / %.1fmm", key, down, up)
		}

	case "remap":
		if len(values) < 1 {
			color.HiRed("Usage: drunkdeer set remap <key> <action>")
			os.Exit(1)
		}
		if state.Remap.Default == nil {
			state.Remap.Default = map[string]string{}
		}
		for _, name := range splitKeys(key) {
			state.Remap.Default[strings.TrimSpace(name)] = values[0]
		}
		rk, ok := driver.GetRemapAction(values[0])
		if !ok {
			color.HiRed("Unknown action: %s", values[0])
			os.Exit(1)
		}
		keys := make(map[int]*driver.RemapKey)
		for _, name := range splitKeys(key) {
			name = strings.TrimSpace(name)
			idx, ok := driver.ResolveKeyToIndex(name, model)
			if !ok {
				color.HiRed("Unknown key: %s", name)
				os.Exit(1)
			}
			keys[idx] = &rk
		}
		a.controller.SendRemapData(keys, 1)
		a.controller.Flush()
		color.HiGreen("Remap %s -> %s", key, values[0])

	default:
		color.HiRed("Unknown section: %s", section)
		color.White("Valid sections: turbo, defaultActuation, light, actuationPoints, rapidTrigger, remap")
		os.Exit(1)
	}

	a.saveSetState(state)
}

func splitKeys(s string) []string {
	return strings.Split(s, ",")
}

func setDefaultBrightness(c *driver.DrunkDeerController) byte {
	if c.Light.Brightness == 0 {
		return 9
	}
	return c.Light.Brightness
}

func sendCustomColor(perKey map[string]string, defaultHex string, c *driver.DrunkDeerController) {
	rgb := parseHexColor(defaultHex)
	c.Light.DefaultColor = rgb
	br := setDefaultBrightness(c)
	cmap := make(map[int][3]byte)
	model := c.GetIdentity().KeyboardModel
	for name, hex := range perKey {
		idx, ok := driver.ResolveKeyToIndex(name, model)
		if ok {
			cmap[idx] = parseHexColor(hex)
		}
	}
	c.SendLEDModeSelect(0x00, driver.SEQUENCE_CUSTOM, c.Light.Speed, br, 0xFF)
	c.Flush()
	time.Sleep(10 * time.Millisecond)
	c.SendCustomColorData(cmap, br, rgb, false)
	c.Flush()
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
	color.HiWhite("  - drunkdeer set <section> <key> [value...]")
	color.HiWhite("  - drunkdeer profiles")
	color.HiWhite("  - drunkdeer reset")
	color.HiWhite("  - drunkdeer list")
	color.HiWhite("  - drunkdeer version")
}
