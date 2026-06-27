package gui

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	canvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	fynetool "fyne.io/fyne/v2/widget"
	"github.com/2xxn/cli-drunkdeer/driver"
	kbwidget "github.com/2xxn/cli-drunkdeer/gui/widget"
)

func udevHint(err error) error {
	return fmt.Errorf("%w\ninstall udev rules: see README at github.com/Sharps1ght/opendrunkdeer", err)
}

func selfExec(args ...string) ([]byte, error) {
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}
	out, err := exec.Command(self, args...).CombinedOutput()
	if err != nil {
		return out, udevHint(err)
	}
	return out, nil
}

func selfExecInteractive(args ...string) error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(self, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return udevHint(err)
	}
	return nil
}

type LightProfile struct {
	Enabled     bool              `json:"enabled"`
	Direction   int               `json:"direction"`
	Speed       int               `json:"speed"`
	Brightness  int               `json:"brightness"`
	Sequence    int               `json:"sequence"`
	ColorIndex  int               `json:"colorIndex,omitempty"`
	Color       string            `json:"color"`
	Colors     map[string]string `json:"colors,omitempty"`
	TurboColor string            `json:"colorTurbo,omitempty"`
}

type RapidTriggerProfile struct {
	Enabled           bool    `json:"enabled"`
	DefaultDownstroke float32 `json:"defaultDownstroke"`
	DefaultUpstroke   float32 `json:"defaultUpstroke"`
}

type RemapProfile struct {
	Default map[string]string `json:"Default,omitempty"`
	Fn      map[string]string `json:"Fn,omitempty"`
	Menu    map[string]string `json:"Menu,omitempty"`
}

type Profile struct {
	Model            string              `json:"model"`
	Turbo            bool                `json:"turbo"`
	DefaultActuation float32             `json:"defaultActuation"`
	RapidTrigger     RapidTriggerProfile `json:"rapidTrigger"`
	Light            LightProfile        `json:"light"`
	ActuationPoints  map[string]float32  `json:"actuationPoints,omitempty"`
	RapidTriggers    map[string][2]float32 `json:"rapidTriggers,omitempty"`
	Remap            RemapProfile        `json:"remap,omitempty"`
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}
	return filepath.Join(home, ".config", "drunkdeer")
}

func loadProfile(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func saveProfile(path string, p *Profile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func scanProfiles(dir string) []string {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && strings.HasSuffix(name, ".json") {
			base := strings.TrimSuffix(name, ".json")
			if strings.HasPrefix(base, ".") {
				continue
			}
			names = append(names, base)
		}
	}
	sort.Strings(names)
	return names
}

func defaultProfile(model string) *Profile {
	return &Profile{
		Model:            model,
		Turbo:            false,
		DefaultActuation: 2.0,
		RapidTrigger: RapidTriggerProfile{
			Enabled:           false,
			DefaultDownstroke: 2.0,
			DefaultUpstroke:   2.0,
		},
		Light: LightProfile{
			Enabled:    true,
			Direction:  0,
			Speed:      9,
			Brightness: 9,
			Sequence:   19,
			Color:      "#00FF00",
			Colors:     map[string]string{},
		},
		ActuationPoints: map[string]float32{},
		RapidTriggers:   map[string][2]float32{},
		Remap:           RemapProfile{},
	}
}

func resolveColorMap(colors map[string]string, model string) map[int]string {
	result := make(map[int]string)
	for name, hex := range colors {
		idx := driver.GetIndexByKey(name, model)
		if idx >= 0 {
			result[idx] = hex
		}
	}
	return result
}

var sequenceColorType = map[int]int{
	0:  0, // Off
	1:  0, // Rotate Marquee
	2:  1, // Wave
	3:  0, // Surf Right
	4:  1, // Breath
	5:  0, // Surf Center
	6:  0, // Cycle
	7:  1, // Ripple
	8:  2, // Always On
	9:  1, // Press
	10: 0, // Snake
	11: 0, // Fountain
	12: 1, // Laser
	13: 0, // Fish
	14: 0, // Surf Cross
	15: 0, // Heart
	16: 0, // Traffic
	17: 0, // Snake Game
	18: 1, // Raindrop
	19: 0, // Custom Light
}

var colorOptions = []struct {
	Name        string
	FirmwareIdx byte
	DisplayHex  string
}{
	{"Rainbow", 0, ""},
	{"Red", 1, "#FF0000"},
	{"Green", 2, "#00FF00"},
	{"Blue", 3, "#0000FF"},
	{"Yellow", 4, "#FFFF00"},
	{"Magenta", 5, "#FF00FF"},
	{"Cyan", 6, "#00FFFF"},
	{"White", 7, "#FFFFFF"},
}

func buildRainbowColors(model string) map[int]string {
	layout := kbwidget.GetLayoutDef(model)
	colors := make(map[int]string)
	for ri, row := range layout.Rows {
		for ci, key := range row {
			hue := float64(ri*len(layout.Rows[0]) + ci) * 15.0
			for hue >= 360 {
				hue -= 360
			}
			r, g, b := hsvToRGB(hue, 1.0, 1.0)
			colors[key.Value] = fmt.Sprintf("#%02x%02x%02x", r, g, b)
		}
	}
	return colors
}

func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	h = float64(int(h) % 360)
	c := v * s
	x := c * (1 - absFloat64(float64(int(h/60)%2)-1))
	m := v - c
	var r, g, b float64
	switch int(h / 60) {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3:
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	case 5:
		r, g, b = c, 0, x
	}
	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}

func absFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func applyKeyboardColors(kb *kbwidget.KeyboardWidget, profile *Profile, model string) {
	if profile.Light.Sequence != 19 {
		ct := sequenceColorType[profile.Light.Sequence]
		if ct == 0 || (ct == 1 && profile.Light.ColorIndex == 0) {
			kb.SetColors("", buildRainbowColors(model))
			return
		}
		if ct > 0 {
			for _, opt := range colorOptions {
				if opt.FirmwareIdx == byte(profile.Light.ColorIndex) {
					if opt.DisplayHex != "" {
						kb.SetColors(opt.DisplayHex, nil)
					} else {
						kb.SetColors("", buildRainbowColors(model))
					}
					return
				}
			}
		}
		kb.SetColors("", nil)
		return
	}
	if profile.Turbo {
		tc := profile.Light.TurboColor
		if tc == "" {
			tc = profile.Light.Color
		}
		kb.SetColors(tc, nil)
	} else {
		kb.SetColors(profile.Light.Color, resolveColorMap(profile.Light.Colors, model))
	}
}

func applyProfileToKeyboard(profile *Profile, model string, dir string) {
	tmpPath := filepath.Join(dir, ".gui-apply.json")
	if err := saveProfile(tmpPath, profile); err != nil {
		log.Printf("Error writing temp profile: %v", err)
		return
	}
	if err := selfExecInteractive("--debug", "load", tmpPath); err != nil {
		log.Printf("Error applying profile: %v", err)
	}
	os.Remove(tmpPath)
}

type floatSliderOpts struct {
	label     string
	initial   float64
	min, max  float64
	step      float64
	onChanged func(v float64)
}

func newFloatSlider(opts floatSliderOpts) (sl *fynetool.Slider, entry *fynetool.Entry, box *fyne.Container) {
	l := fynetool.NewLabel(opts.label)
	l.Alignment = fyne.TextAlignCenter

	sl = fynetool.NewSlider(opts.min, opts.max)
	sl.Step = opts.step
	sl.Value = opts.initial
	sl.Orientation = fynetool.Vertical

	entry = fynetool.NewEntry()
	entry.SetText(fmt.Sprintf("%.1fmm", opts.initial))

	var updating bool
	sl.OnChanged = func(v float64) {
		if updating {
			return
		}
		updating = true
		entry.SetText(fmt.Sprintf("%.1fmm", v))
		if opts.onChanged != nil {
			opts.onChanged(v)
		}
		updating = false
	}

	entry.OnChanged = func(s string) {
		if updating {
			return
		}
		s = strings.TrimSuffix(s, "mm")
		s = strings.TrimSpace(s)
		v, err := strconv.ParseFloat(s, 64)
		if err != nil || v < opts.min || v > opts.max {
			return
		}
		updating = true
		sl.Value = v
		sl.Refresh()
		if opts.onChanged != nil {
			opts.onChanged(v)
		}
		updating = false
	}

	wrapped := kbwidget.NewMinSizeWrap(sl, fyne.NewSize(80, 0))
	box = container.NewBorder(l, entry, nil, nil, wrapped)
	return
}

type intSliderOpts struct {
	label     string
	initial   int
	min, max  int
	onChanged func(v int)
}

func newIntSlider(opts intSliderOpts) (sl *fynetool.Slider, entry *fynetool.Entry, box *fyne.Container) {
	l := fynetool.NewLabel(opts.label)
	l.Alignment = fyne.TextAlignCenter

	sl = fynetool.NewSlider(float64(opts.min), float64(opts.max))
	sl.Step = 1
	sl.Value = float64(opts.initial)
	sl.Orientation = fynetool.Vertical

	entry = fynetool.NewEntry()
	entry.SetText(strconv.Itoa(opts.initial))

	var updating bool
	sl.OnChanged = func(v float64) {
		if updating {
			return
		}
		updating = true
		iv := int(math.Round(v))
		entry.SetText(strconv.Itoa(iv))
		if opts.onChanged != nil {
			opts.onChanged(iv)
		}
		updating = false
	}

	entry.OnChanged = func(s string) {
		if updating {
			return
		}
		s = strings.TrimSpace(s)
		v, err := strconv.Atoi(s)
		if err != nil || v < opts.min || v > opts.max {
			return
		}
		updating = true
		sl.Value = float64(v)
		sl.Refresh()
		if opts.onChanged != nil {
			opts.onChanged(v)
		}
		updating = false
	}

	wrapped := kbwidget.NewMinSizeWrap(sl, fyne.NewSize(80, 0))
	box = container.NewBorder(l, entry, nil, nil, wrapped)
	return
}

func Run() {
	checkDisplay()

	a := app.NewWithID("drunkdeer-config")

	if data := kbwidget.MonospaceFontData(); len(data) > 0 {
		monoFont = fyne.NewStaticResource("Monospace.ttf", data)
		a.Settings().SetTheme(&monoTheme{})
	}

	w := a.NewWindow("DrunkDeer Config")

	modelFlag := flag.String("model", driver.KEYBOARD_G60, "keyboard model (G60, G65, G75, A75)")
	flag.Parse()
	model := *modelFlag
	switch model {
	case driver.KEYBOARD_G60, driver.KEYBOARD_G65, driver.KEYBOARD_G75, driver.KEYBOARD_A75:
	default:
		log.Printf("Unknown model %s, falling back to %s", model, driver.KEYBOARD_G60)
		model = driver.KEYBOARD_G60
	}
	log.Printf("Using keyboard model: %s", model)

	dir := configDir()
	os.MkdirAll(dir, 0755)

	profileNames := scanProfiles(dir)
	if len(profileNames) == 0 {
		profileNames = []string{"default"}
	}

	currentProfile := profileNames[0]
	profilePath := filepath.Join(dir, currentProfile+".json")
	profile, err := loadProfile(profilePath)
	if err != nil {
		log.Printf("No existing profile, using defaults: %v", err)
		profile = defaultProfile(model)
	}
	initialModel := profile.Model

	kb := kbwidget.NewKeyboardWidget(model)
	applyKeyboardColors(kb, profile, model)

	actuationSlider, actuationValue, actuationBox := newFloatSlider(floatSliderOpts{
		label:     "Key\nDown",
		initial:   float64(profile.DefaultActuation),
		min:       0.2,
		max:       3.8,
		step:      0.1,
		onChanged: func(v float64) { profile.DefaultActuation = float32(v) },
	})

	rtDownSlider, rtDownValue, rtDownBox := newFloatSlider(floatSliderOpts{
		label:   "RT\nDown",
		initial: float64(profile.RapidTrigger.DefaultDownstroke),
		min:     0.2,
		max:     3.8,
		step:    0.1,
		onChanged: func(v float64) {
			sel := kb.SelectedKeys()
			if len(sel) > 0 {
				if profile.RapidTriggers == nil {
					profile.RapidTriggers = make(map[string][2]float32)
				}
				for _, idx := range sel {
					name := driver.GetKeyByIndex(idx, model)
					if name != "" {
						vals := profile.RapidTriggers[name]
						vals[0] = float32(v)
						profile.RapidTriggers[name] = vals
					}
				}
			} else {
				profile.RapidTrigger.DefaultDownstroke = float32(v)
			}
		},
	})

	rtUpSlider, rtUpValue, rtUpBox := newFloatSlider(floatSliderOpts{
		label:   "RT\nUp",
		initial: float64(profile.RapidTrigger.DefaultUpstroke),
		min:     0.2,
		max:     3.8,
		step:    0.1,
		onChanged: func(v float64) {
			sel := kb.SelectedKeys()
			if len(sel) > 0 {
				if profile.RapidTriggers == nil {
					profile.RapidTriggers = make(map[string][2]float32)
				}
				for _, idx := range sel {
					name := driver.GetKeyByIndex(idx, model)
					if name != "" {
						vals := profile.RapidTriggers[name]
						vals[1] = float32(v)
						profile.RapidTriggers[name] = vals
					}
				}
			} else {
				profile.RapidTrigger.DefaultUpstroke = float32(v)
			}
		},
	})

	var brTimer *time.Timer
	brSlider, brValue, brBox := newIntSlider(intSliderOpts{
		label:   "Light\nLevel",
		initial: profile.Light.Brightness + 1,
		min:     1,
		max:     10,
		onChanged: func(v int) {
			profile.Light.Brightness = v - 1
			brTimer = startDebounce(brTimer, 500*time.Millisecond, func() {
				if out, err := selfExec("set", "light", "brightness", strconv.Itoa(v-1)); err != nil {
					log.Printf("Error setting light brightness: %v\n%s", err, out)
				}
			})
		},
	})

	var spTimer *time.Timer
	spSlider, spValue, spBox := newIntSlider(intSliderOpts{
		label:   "Anim\nSpeed",
		initial: profile.Light.Speed + 1,
		min:     1,
		max:     10,
		onChanged: func(v int) {
			profile.Light.Speed = v - 1
			spTimer = startDebounce(spTimer, 500*time.Millisecond, func() {
				if out, err := selfExec("set", "light", "speed", strconv.Itoa(v-1)); err != nil {
					log.Printf("Error setting light speed: %v\n%s", err, out)
				}
			})
		},
	})

	setRTEnabled := func(en bool) {
		if en {
			rtDownSlider.Enable()
			rtDownValue.Enable()
			rtUpSlider.Enable()
			rtUpValue.Enable()
		} else {
			rtDownSlider.Disable()
			rtDownValue.Disable()
			rtUpSlider.Disable()
			rtUpValue.Disable()
		}
	}

	rightPanel := container.NewHBox(actuationBox, rtDownBox, rtUpBox)
	leftPanel := container.NewHBox(brBox, spBox)

	var pickingColor bool
	var colorBtn *fynetool.Button
	var currentPicker *kbwidget.ColorPicker
	var savedSelAtOpen []int
	var savedColorsAtOpen map[string]string
	var pickerPopUp *dismissPopUp
	var remapPopUp *dismissPopUp

	donePicking := func() {
		pickingColor = false
		currentPicker = nil
	}

	cancelColor := func() {
		if pickingColor {
			donePicking()
		}
		profile.Light.Colors = savedColorsAtOpen
		applyKeyboardColors(kb, profile, model)
		kb.SelectKeys(savedSelAtOpen)
		colorBtn.Enable()
		if pickerPopUp != nil {
			pickerPopUp.Hide()
			pickerPopUp = nil
		}
	}

	colorBtn = fynetool.NewButton("Color...", func() {
		if pickingColor {
			if pickerPopUp != nil {
				pickerPopUp.Hide()
				pickerPopUp = nil
			}
			donePicking()
			return
		}

		savedSelAtOpen = kb.SelectedKeys()
		if len(savedSelAtOpen) == 0 {
			return
		}

		savedColorsAtOpen = make(map[string]string)
		for k, v := range profile.Light.Colors {
			savedColorsAtOpen[k] = v
		}

		initialHex := profile.Light.Color
		for k, v := range profile.Light.Colors {
			if driver.GetIndexByKey(k, model) == savedSelAtOpen[0] {
				initialHex = v
				break
			}
		}

		pickingColor = true
		kb.ClearSelection()

		currentPicker = kbwidget.NewColorPicker(initialHex,
			func(hex string) {
				for _, val := range savedSelAtOpen {
					kb.SetIndividualColor(val, hex)
					name := driver.GetKeyByIndex(val, model)
					if name != "" {
						if profile.Light.Colors == nil {
							profile.Light.Colors = make(map[string]string)
						}
						profile.Light.Colors[name] = hex
					}
				}
			},
			func(hex string) {
				pickingColor = false
				kb.SelectKeys(savedSelAtOpen)
				colorBtn.Enable()
				if pickerPopUp != nil {
					pickerPopUp.Hide()
					pickerPopUp = nil
				}
			},
			cancelColor,
		)
		currentPicker.Destroy = func() {
			pickingColor = false
		}

		pickerPopUp = &dismissPopUp{
			PopUp:     fynetool.PopUp{Content: currentPicker.Content, Canvas: w.Canvas()},
			onDismiss: cancelColor,
		}
		pickerPopUp.ExtendBaseWidget(pickerPopUp)
		pickerPopUp.ShowAtRelativePosition(
			fyne.NewPos(colorBtn.Size().Width, 0), colorBtn)
	})
	colorBtn.Disable()

	sequenceValues := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	var sequenceLabels []string
	switch model {
	case driver.KEYBOARD_G75, driver.KEYBOARD_A75:
		sequenceLabels = []string{
			"Off", "Rotate Marquee", "Always On", "Cycle", "Breath",
			"Press", "Stars", "Wave", "Surf Center", "Surf Down",
			"Ripple", "Fish", "Fountain", "Traffic", "Snake Game",
			"Surf Left/Right", "Surf Cross", "Laser", "Random Fountain", "Custom Light",
		}
	default:
		sequenceLabels = []string{
			"Off", "Rotate Marquee", "Wave", "Surf Right", "Breath",
			"Surf Center", "Cycle", "Ripple", "Always On", "Press",
			"Snake", "Fountain", "Laser", "Fish", "Surf Cross",
			"Heart", "Traffic", "Snake Game", "Raindrop", "Custom Light",
		}
	}
	currentSeqIdx := 0
	for i, v := range sequenceValues {
		if v == profile.Light.Sequence {
			currentSeqIdx = i
			break
		}
	}
	var updateColorBtnState func()
	var updateColorSelectState func()
	var remapActive bool

	var startingUp = true
	var updatingColorSelect bool

	sequenceSelect := fynetool.NewSelect(sequenceLabels, func(label string) {
		var seq int
		for i, l := range sequenceLabels {
			if l == label {
				seq = sequenceValues[i]
				profile.Light.Sequence = seq
				break
			}
		}
		applyKeyboardColors(kb, profile, model)
		updateColorBtnState()
		if !startingUp {
			go func() {
				if out, err := selfExec("set", "light", "sequence", strconv.Itoa(seq)); err != nil {
					log.Printf("Error setting light sequence: %v\n%s", err, out)
				}
			}()
		}
		updateColorSelectState()
	})

	var remapBtn *fynetool.Button

	updateColorBtnState = func() {
		if pickingColor {
			return
		}
		sel := len(kb.SelectedKeys()) > 0
		if remapActive {
			colorBtn.Disable()
		} else if profile.Light.Sequence != 19 {
			colorBtn.Disable()
		} else if sel {
			colorBtn.Enable()
		} else {
			colorBtn.Disable()
		}
		if sel {
			if remapBtn != nil {
				remapBtn.Enable()
			}
		} else {
			if remapBtn != nil {
				remapBtn.Disable()
			}
		}
	}
	remapActive = false

	colorLabels := make([]string, len(colorOptions))
	for i, opt := range colorOptions {
		colorLabels[i] = opt.Name
	}
	colorIndexSelect := fynetool.NewSelect(colorLabels, func(label string) {
		for i, l := range colorLabels {
			if l == label {
				profile.Light.ColorIndex = int(colorOptions[i].FirmwareIdx)
				break
			}
		}
		applyKeyboardColors(kb, profile, model)
		if !updatingColorSelect {
			if out, err := selfExec("set", "light", "colorIndex", strconv.Itoa(profile.Light.ColorIndex)); err != nil {
				log.Printf("Error setting light color index: %v\n%s", err, out)
			}
		}
	})

	updateColorSelectState = func() {
		updatingColorSelect = true
		defer func() { updatingColorSelect = false }()
		ct := sequenceColorType[profile.Light.Sequence]
		if ct == 0 {
			colorIndexSelect.Disable()
			return
		}
		var opts []int
		if ct == 2 {
			opts = []int{1, 2, 3, 4, 5, 6, 7} // exclude Rainbow (0)
		} else {
			opts = []int{0, 1, 2, 3, 4, 5, 6, 7}
		}
		optLabels := make([]string, len(opts))
		for i, idx := range opts {
			optLabels[i] = colorLabels[idx]
		}
		colorIndexSelect.Options = optLabels

		sel := ""
		for i, idx := range opts {
			if colorOptions[idx].FirmwareIdx == byte(profile.Light.ColorIndex) {
				sel = optLabels[i]
				break
			}
		}
		if sel == "" {
			sel = optLabels[0]
			profile.Light.ColorIndex = int(colorOptions[opts[0]].FirmwareIdx)
		}
		colorIndexSelect.SetSelected(sel)
		colorIndexSelect.Enable()
		colorIndexSelect.Refresh()
	}

	sequenceSelect.SetSelected(sequenceLabels[currentSeqIdx])

	kb.SetOnSelectionChanged(func() {
		updateColorBtnState()
		sel := kb.SelectedKeys()
		if len(sel) > 0 {
			name := driver.GetKeyByIndex(sel[0], model)
			if name != "" {
				if vals, ok := profile.RapidTriggers[name]; ok {
					rtDownSlider.Value = float64(vals[0])
					rtDownSlider.Refresh()
					rtUpSlider.Value = float64(vals[1])
					rtUpSlider.Refresh()
				}
			}
		}
	})

	apply := fynetool.NewButton("Apply", func() {
		applyProfileToKeyboard(profile, model, dir)
	})
	saveBtn := fynetool.NewButton("Save", func() {
		if err := saveProfile(filepath.Join(dir, currentProfile+".json"), profile); err != nil {
			log.Printf("Error saving profile: %v", err)
		}
	})

	modelMismatchLabel := canvas.NewText("Keyboard model mismatch", color.RGBA{0xff, 0xff, 0x00, 0xff})
	modelMismatchLabel.Hidden = true
	noKeyboardLabel := canvas.NewText("No keyboard detected", color.RGBA{0xff, 0x65, 0x00, 0xff})
	noKeyboardLabel.Hidden = true

	if initialModel != "" && initialModel != model {
		profile = defaultProfile(model)
	}
	closeBtn := fynetool.NewButton("Close", func() {
		w.Close()
	})
	bottomBar := container.NewBorder(nil, nil,
		container.NewHBox(sequenceSelect, colorIndexSelect, colorBtn),
		container.NewHBox(apply, saveBtn, closeBtn),
	)

	profileOptions := append(profileNames, "New...")
	var selecting bool
	var profileSelect *fynetool.Select
	var setMismatchUI func(bool)
	var setNoKeyboardUI func(bool)
	var rtCheck, turboCheck *fynetool.Check
	var turboHexUpdating bool
	var turboTimer *time.Timer
	turboHexEntry := fynetool.NewEntry()

	syncProfileUI := func() {
		applyKeyboardColors(kb, profile, model)
		currentSeqIdx = 0
		for i, v := range sequenceValues {
			if v == profile.Light.Sequence {
				currentSeqIdx = i
				break
			}
		}
		sequenceSelect.SetSelected(sequenceLabels[currentSeqIdx])
		updateColorBtnState()
		updateColorSelectState()
		actuationSlider.Value = float64(profile.DefaultActuation)
		actuationSlider.Refresh()
		actuationValue.SetText(fmt.Sprintf("%.1fmm", profile.DefaultActuation))
		rtDownSlider.Value = float64(profile.RapidTrigger.DefaultDownstroke)
		rtDownSlider.Refresh()
		rtDownValue.SetText(fmt.Sprintf("%.1fmm", profile.RapidTrigger.DefaultDownstroke))
		rtUpSlider.Value = float64(profile.RapidTrigger.DefaultUpstroke)
		rtUpSlider.Refresh()
		rtUpValue.SetText(fmt.Sprintf("%.1fmm", profile.RapidTrigger.DefaultUpstroke))
		brSlider.Value = float64(profile.Light.Brightness + 1)
		brSlider.Refresh()
		brValue.SetText(strconv.Itoa(profile.Light.Brightness + 1))
		spSlider.Value = float64(profile.Light.Speed + 1)
		spSlider.Refresh()
		spValue.SetText(strconv.Itoa(profile.Light.Speed + 1))
		if rtCheck != nil {
			rtCheck.SetChecked(profile.RapidTrigger.Enabled)
		}
		setRTEnabled(profile.RapidTrigger.Enabled)
		if turboCheck != nil {
			turboCheck.SetChecked(profile.Turbo)
		}
		if profile.Turbo {
			turboHexEntry.Enable()
		} else {
			turboHexEntry.Disable()
		}
		turboHexUpdating = true
		if profile.Light.TurboColor != "" {
			turboHexEntry.SetText(profile.Light.TurboColor)
		} else {
			turboHexEntry.SetText(profile.Light.Color)
		}
		turboHexUpdating = false
	}

	profileSelect = fynetool.NewSelect(profileOptions, func(name string) {
		if selecting {
			return
		}
		selecting = true
		defer func() { selecting = false }()

		if name == "New..." {
			nameEntry := fynetool.NewEntry()
			dialog.ShowCustomConfirm("New Profile", "Save", "Cancel",
				container.NewVBox(
					fynetool.NewLabel("Profile name:"),
					nameEntry,
				),
				func(ok bool) {
					if !ok || nameEntry.Text == "" {
						profileSelect.SetSelected(currentProfile)
						return
					}
					name := strings.TrimSpace(nameEntry.Text)
					path := filepath.Join(dir, name+".json")
					if _, err := os.Stat(path); err == nil {
						log.Printf("Profile %s already exists", name)
						profileSelect.SetSelected(currentProfile)
						return
					}
					p := defaultProfile(model)
					if err := saveProfile(path, p); err != nil {
						log.Printf("Error creating profile: %v", err)
						profileSelect.SetSelected(currentProfile)
						return
					}
					newNames := scanProfiles(dir)
					profileSelect.Options = append(newNames, "New...")
					profileSelect.SetSelected(name)
					currentProfile = name

				setMismatchUI(false)
				*profile = *p
				syncProfileUI()
				},
				w,
			)
			return
		}
		currentProfile = name
		path := filepath.Join(dir, name+".json")
		p, err := loadProfile(path)
		if err != nil {
			log.Printf("Error loading profile: %v", err)
			return
		}
		if p.Model != "" && p.Model != model {
			setMismatchUI(true)
			p = defaultProfile(model)
		} else {
			setMismatchUI(false)
		}
		*profile = *p
		syncProfileUI()
	})
	deleteBtn := fynetool.NewButton("Delete", func() {
		if len(profileSelect.Options) <= 2 {
			return
		}
		name := currentProfile
		path := filepath.Join(dir, name+".json")
		os.Remove(path)
		newNames := scanProfiles(dir)
		if len(newNames) == 0 {
			newNames = []string{"default"}
			saveProfile(filepath.Join(dir, "default.json"), defaultProfile(model))
		}
		profileSelect.Options = append(newNames, "New...")
		profileSelect.Refresh()

		currentProfile = ""
		profileSelect.SetSelected(newNames[0])
	})

	rtCheck = fynetool.NewCheck("Rapid Trigger", func(checked bool) {
		profile.RapidTrigger.Enabled = checked
		setRTEnabled(checked)
	})
	rtCheck.SetChecked(profile.RapidTrigger.Enabled)
	setRTEnabled(profile.RapidTrigger.Enabled)

	if profile.Light.TurboColor != "" {
		turboHexEntry.SetText(profile.Light.TurboColor)
	} else {
		turboHexEntry.SetText(profile.Light.Color)
	}
	turboHexEntry.SetPlaceHolder("#RRGGBB")
	turboHexEntry.OnChanged = func(s string) {
		if turboHexUpdating {
			return
		}
		profile.Light.TurboColor = s
		applyKeyboardColors(kb, profile, model)
		turboTimer = startDebounce(turboTimer, 500*time.Millisecond, func() {
			matched, _ := regexp.MatchString("^#[0-9a-fA-F]{6}$", s)
			if matched {
				if out, err := selfExec("set", "light", "colorTurbo", s); err != nil {
					log.Printf("Error setting turbo color: %v\n%s", err, out)
				}
			}
		})
	}
	if profile.Turbo {
		turboHexEntry.Enable()
	} else {
		turboHexEntry.Disable()
	}
	turboHexWrapped := kbwidget.NewMinSizeWrap(turboHexEntry, fyne.NewSize(108, 38))

	turboCheck = fynetool.NewCheck("Turbo", func(checked bool) {
		profile.Turbo = checked
		applyKeyboardColors(kb, profile, model)
		if checked {
			turboHexEntry.Enable()
		} else {
			turboHexEntry.Disable()
		}
		val := "false"
		if checked {
			val = "true"
		}
		if out, err := selfExec("set", "turbo", val); err != nil {
			log.Printf("Error setting turbo: %v\n%s", err, out)
		}
	})
	turboCheck.SetChecked(profile.Turbo)

	buildRemapKeys := func(mode string) (map[int]bool, map[int]string) {
		if mode == "" {
			return nil, nil
		}
		var defaults map[int]bool
		var defActions map[int]string
		var entries map[string]string
		switch mode {
		case "Off":
			entries = profile.Remap.Default
			if len(entries) == 0 {
				return nil, nil
			}
		case "Fn":
			defaults = driver.DefaultFnKeys(model)
			defActions = driver.DefaultFnActions(model)
			entries = profile.Remap.Fn
		case "Menu":
			defaults = driver.DefaultMenuKeys(model)
			defActions = driver.DefaultMenuActions(model)
			entries = profile.Remap.Menu
		}
		keys := make(map[int]bool)
		names := make(map[int]string)
		if defaults != nil {
			for k := range defaults {
				keys[k] = true
				if a, ok := defActions[k]; ok {
					names[k] = driver.ActionDisplayName(a)
				}
			}
		}
		for name, action := range entries {
			idx := driver.GetIndexByKey(name, model)
			if idx < 0 {
				continue
			}
			if action == "" {
				delete(keys, idx)
				delete(names, idx)
				continue
			}
			keys[idx] = true
			names[idx] = driver.ActionDisplayName(action)
		}
		return keys, names
	}

	remapRadio := fynetool.NewRadioGroup([]string{"Off", "Fn", "Menu"}, func(s string) {
		keys, names := buildRemapKeys(s)
		kb.SetRemapMode(s, keys, names)
		remapActive = s != "Off"
		if s == "Off" {
			sequenceSelect.Enable()
			setRTEnabled(profile.RapidTrigger.Enabled)
		} else {
			sequenceSelect.Disable()
			setRTEnabled(false)
		}
		updateColorBtnState()
	})
	remapRadio.Horizontal = true
	remapRadio.Required = true
	remapRadio.SetSelected("Off")

	remapBtn = fynetool.NewButton("Remap...", func() {
		savedSelAtOpen = kb.SelectedKeys()
		if len(savedSelAtOpen) == 0 {
			return
		}

		kb.ClearSelection()

		finishRemap := func() {
			mode := remapRadio.Selected
			keys, names := buildRemapKeys(mode)
			kb.SetRemapMode(mode, keys, names)
			if remapPopUp != nil {
				remapPopUp.Hide()
				remapPopUp = nil
			}
		}
		getEntries := func() map[string]string {
			mode := remapRadio.Selected
			switch mode {
			case "Off":
				if profile.Remap.Default == nil {
					profile.Remap.Default = make(map[string]string)
				}
				return profile.Remap.Default
			case "Fn":
				if profile.Remap.Fn == nil {
					profile.Remap.Fn = make(map[string]string)
				}
				return profile.Remap.Fn
			case "Menu":
				if profile.Remap.Menu == nil {
					profile.Remap.Menu = make(map[string]string)
				}
				return profile.Remap.Menu
			}
			return nil
		}

		remapPicker := kbwidget.NewActionPicker(
			func(action string) {
				entries := getEntries()
				for _, val := range savedSelAtOpen {
					name := driver.GetKeyByIndex(val, model)
					if name == "" {
						continue
					}
					entries[name] = action
				}
				finishRemap()
			},
			func() {
				entries := getEntries()
				for _, val := range savedSelAtOpen {
					name := driver.GetKeyByIndex(val, model)
					if name == "" {
						continue
					}
					delete(entries, name)
				}
				finishRemap()
			},
			func() {
				entries := getEntries()
				for _, val := range savedSelAtOpen {
					name := driver.GetKeyByIndex(val, model)
					if name == "" {
						continue
					}
					entries[name] = ""
				}
				finishRemap()
			},
			func() {
				if remapPopUp != nil {
					remapPopUp.Hide()
					remapPopUp = nil
				}
			},
		)
		remapPicker.Destroy = func() {}

		remapPopUp = &dismissPopUp{
			PopUp:     fynetool.PopUp{Content: remapPicker.Content, Canvas: w.Canvas()},
			onDismiss: func() {},
		}
		remapPopUp.ExtendBaseWidget(remapPopUp)

		canvasSize := w.Canvas().Size()
		popupW := canvasSize.Width * 0.75
		popupH := canvasSize.Height * 0.75
		popupSize := fyne.NewSize(popupW, popupH)
		remapPicker.SetContentMinSize(popupSize)
		remapPopUp.ShowAtPosition(fyne.NewPos(
			(canvasSize.Width-popupSize.Width)/2,
			(canvasSize.Height-popupSize.Height)/2,
		))
	})
	remapBtn.Disable()

	disableWidgets := func(dis bool) {
		widgets := []fyne.Disableable{
			profileSelect, apply, saveBtn, deleteBtn,
			colorBtn, colorIndexSelect,
			actuationSlider, actuationValue,
			rtDownSlider, rtDownValue, rtUpSlider, rtUpValue,
			rtCheck, turboCheck, turboHexEntry,
			brSlider, brValue, spSlider, spValue,
		}
		for _, w := range widgets {
			if dis {
				w.Disable()
			} else {
				w.Enable()
			}
		}
		if dis {
			remapBtn.Disable()
		} else {
			updateColorBtnState()
			updateColorSelectState()
		}
	}
	setMismatchUI = func(mismatch bool) {
		modelMismatchLabel.Hidden = !mismatch
		disableWidgets(mismatch)
	}
	setNoKeyboardUI = func(nokey bool) {
		noKeyboardLabel.Hidden = !nokey
		if nokey {
			modelMismatchLabel.Hidden = true
		}
		disableWidgets(!nokey)
	}
	if initialModel != "" && initialModel != model {
		setMismatchUI(true)
	}
	profileSelect.SetSelected(currentProfile)
	topBar := container.NewBorder(nil, nil, container.NewHBox(profileSelect, deleteBtn, rtCheck, turboCheck, turboHexWrapped, remapRadio, remapBtn), container.NewVBox(modelMismatchLabel, noKeyboardLabel))

	kbCentered := container.NewCenter(
		container.NewVBox(
			layout.NewSpacer(),
			kb,
			layout.NewSpacer(),
		),
	)

	content := container.NewBorder(topBar, bottomBar, leftPanel, rightPanel, kbCentered)
	w.SetContent(content)

	pad := float32(96)
	butH := bottomBar.MinSize().Height
	panelW := float32(360)
	updateColorSelectState()
	if !profile.Turbo {
		turboHexEntry.Disable()
	}
	startingUp = false
	w.Resize(fyne.NewSize(
		kb.MinSize().Width+pad+panelW,
		kb.MinSize().Height+pad+butH+float32(60),
	))
	w.CenterOnScreen()
	out, err := selfExec("--list")
	if err != nil || strings.Contains(string(out), "No devices found") {
		model = driver.KEYBOARD_A75
		kb.SetModel(model)
		profile = defaultProfile(model)
		setNoKeyboardUI(true)
	}
	w.ShowAndRun()
}

func startDebounce(prev *time.Timer, delay time.Duration, fn func()) *time.Timer {
	if prev != nil {
		prev.Stop()
	}
	return time.AfterFunc(delay, fn)
}

var monoFont fyne.Resource

type monoTheme struct{}

type dismissPopUp struct {
	fynetool.PopUp
	onDismiss func()
}

func (p *dismissPopUp) Tapped(e *fyne.PointEvent) {
	cPos := p.Content.Position()
	cSize := p.Content.Size()
	if e.Position.X < cPos.X || e.Position.X > cPos.X+cSize.Width ||
		e.Position.Y < cPos.Y || e.Position.Y > cPos.Y+cSize.Height {
		if p.onDismiss != nil {
			p.onDismiss()
		}
		p.Hide()
		return
	}
}

func (p *dismissPopUp) TappedSecondary(e *fyne.PointEvent) {
	p.Tapped(e)
}

func (m *monoTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}
func (m *monoTheme) Font(style fyne.TextStyle) fyne.Resource {
	return monoFont
}
func (m *monoTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (m *monoTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n) * 1.2
}
