package main

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

type LightProfile struct {
	Enabled     bool              `json:"enabled"`
	Direction   int               `json:"direction"`
	Speed       int               `json:"speed"`
	Brightness  int               `json:"brightness"`
	Sequence    int               `json:"sequence"`
	Color       string            `json:"color"`
	Colors     map[string]string `json:"colors,omitempty"`
	TurboColor string            `json:"turboColor,omitempty"`
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
	ActuationPoints  map[string]float32  `json:"actuationPoints"`
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

func applyKeyboardColors(kb *kbwidget.KeyboardWidget, profile *Profile, model string) {
	if profile.Light.Sequence != 19 {
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
	cmd := exec.Command("sudo", "-E", "drunkdeer-cli", "--debug", "load", tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error applying profile: %v", err)
	}
	os.Remove(tmpPath)
}

func main() {
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		fmt.Fprintln(os.Stderr, "No display server found (set DISPLAY for X11 or WAYLAND_DISPLAY for Wayland)")
		os.Exit(1)
	}

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

	var actUpdating bool

	actuationLabel := fynetool.NewLabel("Actuation\nPoint")
	actuationLabel.Alignment = fyne.TextAlignCenter

	actuationSlider := fynetool.NewSlider(0.2, 3.8)
	actuationSlider.Step = 0.1
	actuationSlider.Value = float64(profile.DefaultActuation)
	actuationSlider.Orientation = fynetool.Vertical

	actuationValue := fynetool.NewEntry()
	actuationValue.SetText(fmt.Sprintf("%.1fmm", profile.DefaultActuation))

	actuationSlider.OnChanged = func(v float64) {
		if actUpdating {
			return
		}
		actUpdating = true
		actuationValue.SetText(fmt.Sprintf("%.1fmm", v))
		profile.DefaultActuation = float32(v)
		actUpdating = false
	}

	actuationValue.OnChanged = func(s string) {
		if actUpdating {
			return
		}
		s = strings.TrimSuffix(s, "mm")
		s = strings.TrimSpace(s)
		v, err := strconv.ParseFloat(s, 64)
		if err != nil || v < 0.2 || v > 3.8 {
			return
		}
		actUpdating = true
		actuationSlider.Value = v
		actuationSlider.Refresh()
		profile.DefaultActuation = float32(v)
		actUpdating = false
	}

	actuationBox := container.NewBorder(
		actuationLabel,
		actuationValue,
		nil, nil,
		actuationSlider,
	)

	var rtUpdating bool

	rtLabel := fynetool.NewLabel("Rapid\nTrigger")
	rtLabel.Alignment = fyne.TextAlignCenter

	rtSlider := fynetool.NewSlider(0.2, 3.8)
	rtSlider.Step = 0.1
	rtSlider.Value = float64(profile.RapidTrigger.DefaultDownstroke)
	rtSlider.Orientation = fynetool.Vertical

	rtValue := fynetool.NewEntry()
	rtValue.SetText(fmt.Sprintf("%.1fmm", profile.RapidTrigger.DefaultDownstroke))

	rtSlider.OnChanged = func(v float64) {
		if rtUpdating {
			return
		}
		rtUpdating = true
		rtValue.SetText(fmt.Sprintf("%.1fmm", v))
		profile.RapidTrigger.DefaultDownstroke = float32(v)
		profile.RapidTrigger.DefaultUpstroke = float32(v)
		rtUpdating = false
	}

	rtValue.OnChanged = func(s string) {
		if rtUpdating {
			return
		}
		s = strings.TrimSuffix(s, "mm")
		s = strings.TrimSpace(s)
		v, err := strconv.ParseFloat(s, 64)
		if err != nil || v < 0.2 || v > 3.8 {
			return
		}
		rtUpdating = true
		rtSlider.Value = v
		rtSlider.Refresh()
		profile.RapidTrigger.DefaultDownstroke = float32(v)
		profile.RapidTrigger.DefaultUpstroke = float32(v)
		rtUpdating = false
	}

	rtBox := container.NewBorder(
		rtLabel,
		rtValue,
		nil, nil,
		rtSlider,
	)

	var brUpdating bool
	var brTimer *time.Timer

	brLabel := fynetool.NewLabel("Brightness\nLevel")
	brLabel.Alignment = fyne.TextAlignCenter

	brSlider := fynetool.NewSlider(1, 10)
	brSlider.Step = 1
	brSlider.Value = float64(profile.Light.Brightness + 1)
	brSlider.Orientation = fynetool.Vertical

	brValue := fynetool.NewEntry()
	brValue.SetText(strconv.Itoa(profile.Light.Brightness + 1))

	brSlider.OnChanged = func(v float64) {
		if brUpdating {
			return
		}
		brUpdating = true
		iv := int(math.Round(v))
		brValue.SetText(strconv.Itoa(iv))
		profile.Light.Brightness = iv - 1
		brUpdating = false
		if brTimer != nil {
			brTimer.Stop()
		}
		brTimer = time.AfterFunc(500*time.Millisecond, func() {
			cmd := exec.Command("sudo", "-E", "drunkdeer-cli", "set", "light", "brightness", strconv.Itoa(iv-1))
			if out, err := cmd.CombinedOutput(); err != nil {
				log.Printf("Error setting light brightness: %v\n%s", err, out)
			}
		})
	}

	brValue.OnChanged = func(s string) {
		if brUpdating {
			return
		}
		s = strings.TrimSpace(s)
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 || v > 10 {
			return
		}
		brUpdating = true
		brSlider.Value = float64(v)
		brSlider.Refresh()
		profile.Light.Brightness = v - 1
		brUpdating = false
		if brTimer != nil {
			brTimer.Stop()
		}
		brTimer = time.AfterFunc(500*time.Millisecond, func() {
			cmd := exec.Command("sudo", "-E", "drunkdeer-cli", "set", "light", "brightness", strconv.Itoa(v-1))
			if out, err := cmd.CombinedOutput(); err != nil {
				log.Printf("Error setting light brightness: %v\n%s", err, out)
			}
		})
	}

	brBox := container.NewBorder(
		brLabel,
		brValue,
		nil, nil,
		brSlider,
	)

	var spUpdating bool
	var spTimer *time.Timer

	spLabel := fynetool.NewLabel("Animation\nSpeed")
	spLabel.Alignment = fyne.TextAlignCenter

	spSlider := fynetool.NewSlider(1, 10)
	spSlider.Step = 1
	spSlider.Value = float64(profile.Light.Speed + 1)
	spSlider.Orientation = fynetool.Vertical

	spValue := fynetool.NewEntry()
	spValue.SetText(strconv.Itoa(profile.Light.Speed + 1))

	spSlider.OnChanged = func(v float64) {
		if spUpdating {
			return
		}
		spUpdating = true
		iv := int(math.Round(v))
		spValue.SetText(strconv.Itoa(iv))
		profile.Light.Speed = iv - 1
		spUpdating = false
		if spTimer != nil {
			spTimer.Stop()
		}
		spTimer = time.AfterFunc(500*time.Millisecond, func() {
			cmd := exec.Command("sudo", "-E", "drunkdeer-cli", "set", "light", "speed", strconv.Itoa(iv-1))
			if out, err := cmd.CombinedOutput(); err != nil {
				log.Printf("Error setting light speed: %v\n%s", err, out)
			}
		})
	}

	spValue.OnChanged = func(s string) {
		if spUpdating {
			return
		}
		s = strings.TrimSpace(s)
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 || v > 10 {
			return
		}
		spUpdating = true
		spSlider.Value = float64(v)
		spSlider.Refresh()
		profile.Light.Speed = v - 1
		spUpdating = false
		if spTimer != nil {
			spTimer.Stop()
		}
		spTimer = time.AfterFunc(500*time.Millisecond, func() {
			cmd := exec.Command("sudo", "-E", "drunkdeer-cli", "set", "light", "speed", strconv.Itoa(v-1))
			if out, err := cmd.CombinedOutput(); err != nil {
				log.Printf("Error setting light speed: %v\n%s", err, out)
			}
		})
	}

	spBox := container.NewBorder(
		spLabel,
		spValue,
		nil, nil,
		spSlider,
	)

	rightPanel := container.NewHBox(actuationBox, rtBox)
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
	var remapActive bool

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
		go func() {
			cmd := exec.Command("sudo", "-E", "drunkdeer-cli", "set", "light", "sequence", strconv.Itoa(seq))
			if out, err := cmd.CombinedOutput(); err != nil {
				log.Printf("Error setting light sequence: %v\n%s", err, out)
			}
		}()
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
	sequenceSelect.SetSelected(sequenceLabels[currentSeqIdx])

	kb.SetOnSelectionChanged(func() {
		updateColorBtnState()
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

	if initialModel != "" && initialModel != model {
		profile = defaultProfile(model)
	}
	closeBtn := fynetool.NewButton("Close", func() {
		w.Close()
	})
	bottomBar := container.NewBorder(nil, nil,
		container.NewHBox(sequenceSelect, colorBtn),
		container.NewHBox(apply, saveBtn, closeBtn),
	)

	profileOptions := append(profileNames, "New...")
	var selecting bool
	var profileSelect *fynetool.Select
	var setMismatchUI func(bool)
	var rtCheck, turboCheck *fynetool.Check
	var turboHexUpdating bool
	turboHexEntry := fynetool.NewEntry()
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
					profile.DefaultActuation = p.DefaultActuation
					profile.RapidTrigger = p.RapidTrigger
					profile.Turbo = p.Turbo
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
					actuationSlider.Value = float64(p.DefaultActuation)
					actuationSlider.Refresh()
					actuationValue.SetText(fmt.Sprintf("%.1fmm", p.DefaultActuation))
					rtSlider.Value = float64(p.RapidTrigger.DefaultDownstroke)
					rtSlider.Refresh()
					rtValue.SetText(fmt.Sprintf("%.1fmm", p.RapidTrigger.DefaultDownstroke))
					brSlider.Value = float64(p.Light.Brightness + 1)
					brSlider.Refresh()
					brValue.SetText(strconv.Itoa(p.Light.Brightness + 1))
					spSlider.Value = float64(p.Light.Speed + 1)
					spSlider.Refresh()
					spValue.SetText(strconv.Itoa(p.Light.Speed + 1))
					if rtCheck != nil {
						rtCheck.SetChecked(p.RapidTrigger.Enabled)
						if p.RapidTrigger.Enabled {
							rtSlider.Enable()
							rtValue.Enable()
						} else {
							rtSlider.Disable()
							rtValue.Disable()
						}
					}
					if turboCheck != nil {
						turboCheck.SetChecked(p.Turbo)
					}
					turboHexUpdating = true
					if p.Light.TurboColor != "" {
						turboHexEntry.SetText(p.Light.TurboColor)
					} else {
						turboHexEntry.SetText(p.Light.Color)
					}
					turboHexUpdating = false
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
		actuationSlider.Value = float64(profile.DefaultActuation)
		actuationSlider.Refresh()
		actuationValue.SetText(fmt.Sprintf("%.1fmm", profile.DefaultActuation))
		rtSlider.Value = float64(profile.RapidTrigger.DefaultDownstroke)
		rtSlider.Refresh()
		rtValue.SetText(fmt.Sprintf("%.1fmm", profile.RapidTrigger.DefaultDownstroke))
		brSlider.Value = float64(profile.Light.Brightness + 1)
		brSlider.Refresh()
		brValue.SetText(strconv.Itoa(profile.Light.Brightness + 1))
		spSlider.Value = float64(profile.Light.Speed + 1)
		spSlider.Refresh()
		spValue.SetText(strconv.Itoa(profile.Light.Speed + 1))
		if rtCheck != nil {
			rtCheck.SetChecked(profile.RapidTrigger.Enabled)
			if profile.RapidTrigger.Enabled {
				rtSlider.Enable()
				rtValue.Enable()
			} else {
				rtSlider.Disable()
				rtValue.Disable()
			}
		}
		if turboCheck != nil {
			turboCheck.SetChecked(profile.Turbo)
		}
		turboHexUpdating = true
		if profile.Light.TurboColor != "" {
			turboHexEntry.SetText(profile.Light.TurboColor)
		} else {
			turboHexEntry.SetText(profile.Light.Color)
		}
		turboHexUpdating = false
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
		if checked {
			rtSlider.Enable()
			rtValue.Enable()
		} else {
			rtSlider.Disable()
			rtValue.Disable()
		}
	})
	rtCheck.SetChecked(profile.RapidTrigger.Enabled)
	if !profile.RapidTrigger.Enabled {
		rtSlider.Disable()
		rtValue.Disable()
	}

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
	}
	if profile.Turbo {
		turboHexEntry.Enable()
	} else {
		turboHexEntry.Disable()
	}
	turboHexWrapped := &minSizeWrap{inner: turboHexEntry, minsize: fyne.NewSize(108, 38)}
	turboHexWrapped.ExtendBaseWidget(turboHexWrapped)

	turboCheck = fynetool.NewCheck("Turbo", func(checked bool) {
		profile.Turbo = checked
		applyKeyboardColors(kb, profile, model)
		if checked {
			turboHexEntry.Enable()
		} else {
			turboHexEntry.Disable()
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
		} else {
			sequenceSelect.Disable()
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

	setMismatchUI = func(mismatch bool) {
		modelMismatchLabel.Hidden = !mismatch
		if mismatch {
			profileSelect.Disable()
			apply.Disable()
			saveBtn.Disable()
			deleteBtn.Disable()
			colorBtn.Disable()
			actuationSlider.Disable()
			actuationValue.Disable()
			rtSlider.Disable()
			rtValue.Disable()
			rtCheck.Disable()
			turboCheck.Disable()
			turboHexEntry.Disable()
			brSlider.Disable()
			brValue.Disable()
			spSlider.Disable()
			spValue.Disable()
			remapBtn.Disable()
		} else {
			profileSelect.Enable()
			apply.Enable()
			saveBtn.Enable()
			deleteBtn.Enable()
			updateColorBtnState()
			actuationSlider.Enable()
			actuationValue.Enable()
			rtSlider.Enable()
			rtValue.Enable()
			rtCheck.Enable()
			turboCheck.Enable()
			turboHexEntry.Enable()
			brSlider.Enable()
			brValue.Enable()
			spSlider.Enable()
			spValue.Enable()
		}
	}
	if initialModel != "" && initialModel != model {
		setMismatchUI(true)
	}
	profileSelect.SetSelected(currentProfile)
	topBar := container.NewBorder(nil, nil, container.NewHBox(profileSelect, deleteBtn, rtCheck, turboCheck, turboHexWrapped, remapRadio, remapBtn), modelMismatchLabel)

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
	w.Resize(fyne.NewSize(
		kb.MinSize().Width+pad+panelW,
		kb.MinSize().Height+pad+butH+float32(60),
	))
	w.CenterOnScreen()
	w.ShowAndRun()
}

type minSizeWrap struct {
	fynetool.BaseWidget
	inner   fyne.CanvasObject
	minsize fyne.Size
}

func (m *minSizeWrap) CreateRenderer() fyne.WidgetRenderer {
	return &minSizeRender{obj: m.inner, minSize: m.minsize}
}

type minSizeRender struct {
	obj     fyne.CanvasObject
	minSize fyne.Size
}

func (r *minSizeRender) Objects() []fyne.CanvasObject              { return []fyne.CanvasObject{r.obj} }
func (r *minSizeRender) Layout(s fyne.Size)                          { r.obj.Resize(s) }
func (r *minSizeRender) MinSize() fyne.Size                          { return r.minSize }
func (r *minSizeRender) Refresh()                                    {}
func (r *minSizeRender) ApplyTheme()                                 {}
func (r *minSizeRender) BackgroundColor() color.Color                { return color.Transparent }
func (r *minSizeRender) Destroy()                                    {}

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
