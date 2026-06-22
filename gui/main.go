package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

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

	rightPanel := container.NewHBox(
		actuationBox,
		rtBox,
	)

	var pickingColor bool
	var colorBtn *fynetool.Button
	var currentPicker *kbwidget.ColorPicker
	var savedSelAtOpen []int
	var savedColorsAtOpen map[string]string
	pickerContainer := container.New(layout.NewMaxLayout())

	donePicking := func() {
		pickingColor = false
		pickerContainer.RemoveAll()
		currentPicker = nil
	}

	cancelColor := func() {
		donePicking()
		profile.Light.Colors = savedColorsAtOpen
		applyKeyboardColors(kb, profile, model)
		kb.SelectKeys(savedSelAtOpen)
		colorBtn.Enable()
	}

	colorBtn = fynetool.NewButton("Color...", func() {
		if pickingColor {
			cancelColor()
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
				donePicking()
				kb.SelectKeys(savedSelAtOpen)
				colorBtn.Enable()
			},
			cancelColor,
		)
		currentPicker.Destroy = donePicking
		pickerContainer.RemoveAll()
		pickerContainer.Add(currentPicker.Content)
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

	updateColorBtnState = func() {
		if pickingColor || remapActive || profile.Light.Sequence != 19 {
			colorBtn.Disable()
			return
		}
		if len(kb.SelectedKeys()) > 0 {
			colorBtn.Enable()
		} else {
			colorBtn.Disable()
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
		pickerContainer,
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
					if rtCheck != nil {
						rtCheck.SetChecked(p.RapidTrigger.Enabled)
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
		if rtCheck != nil {
			rtCheck.SetChecked(profile.RapidTrigger.Enabled)
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
	})
	rtCheck.SetChecked(profile.RapidTrigger.Enabled)

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
	turboHexWrapped := &minSizeWrap{inner: turboHexEntry, minsize: fyne.NewSize(90, 32)}
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
			if idx >= 0 {
				keys[idx] = true
				names[idx] = driver.ActionDisplayName(action)
			}
		}
		return keys, names
	}

	remapRadio := fynetool.NewRadioGroup([]string{"Off", "Fn", "Menu"}, func(s string) {
		if s == "Off" {
			kb.SetRemapMode("", nil, nil)
			remapActive = false
		} else {
			keys, names := buildRemapKeys(s)
			kb.SetRemapMode(s, keys, names)
			remapActive = true
		}
		updateColorBtnState()
	})
	remapRadio.Horizontal = true
	remapRadio.Required = true
	remapRadio.SetSelected("Off")

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
		}
	}
	if initialModel != "" && initialModel != model {
		setMismatchUI(true)
	}
	profileSelect.SetSelected(currentProfile)
	topBar := container.NewBorder(nil, nil, container.NewHBox(profileSelect, deleteBtn, rtCheck, turboCheck, turboHexWrapped, remapRadio), modelMismatchLabel)

	kbCentered := container.NewCenter(
		container.NewVBox(
			layout.NewSpacer(),
			kb,
			layout.NewSpacer(),
		),
	)

	content := container.NewBorder(topBar, bottomBar, nil, rightPanel, kbCentered)
	w.SetContent(content)

	pad := float32(80)
	butH := bottomBar.MinSize().Height
	panelW := float32(300)
	w.Resize(fyne.NewSize(
		kb.MinSize().Width+pad+panelW,
		kb.MinSize().Height+pad+butH+float32(50),
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
	return theme.DefaultTheme().Size(n)
}
