package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetool "fyne.io/fyne/v2/widget"
	"github.com/2xxn/cli-drunkdeer/driver"
)

type ActionPicker struct {
	Content fyne.CanvasObject
	Destroy func()
	scroll  *container.Scroll
}

func (a *ActionPicker) SetContentMinSize(s fyne.Size) {
	if a.scroll != nil {
		a.scroll.SetMinSize(s)
	}
}

type actionGroup struct {
	Title   string
	Actions []string
}

var actionGroups = []actionGroup{
	{
		"Features",
		[]string{
			driver.KeyLightModeSw, driver.KeyLightModeCycle, driver.KeyLightColorCycle,
			driver.KeyLightModeNext, driver.KeyLightModePrev, driver.KeyLightPause,
			driver.KeyLightLumiInc, driver.KeyLightLumiDec,
			driver.KeyLightSpeedInc, driver.KeyLightSpeedDec,
			driver.KeyLightCtrl, driver.KeyColorCtrl, driver.KeySpeedCtrl, driver.KeyBrCtrl,
			driver.KeyWinLock,
			driver.KeyRdt, driver.KeyLw,
			driver.KeyRt, driver.KeyRtMatch, driver.KeyRtExtreme, driver.KeyRtStandard, driver.KeyRtCompete,
			driver.KeyTurbo,
		},
	},
	{
		"Characters",
		[]string{
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"MINUS", "PLUS", "BRKTS_L", "BRKTS_R", "COLON", "QOTATN",
			"TILDE", "COMMA", "PERIOD", "SLASH", "CAPS",
			driver.KeyKana, driver.KeyLoop, driver.KeyNoLoop,
		},
	},
	{
		"Numbers",
		[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
	},
	{
		"Numpad",
		[]string{
			driver.KeyNum0, driver.KeyNum1, driver.KeyNum2, driver.KeyNum3, driver.KeyNum4,
			driver.KeyNum5, driver.KeyNum6, driver.KeyNum7, driver.KeyNum8, driver.KeyNum9,
			driver.KeyKpPlus, driver.KeyKpMinus, driver.KeyKpMult, driver.KeyKpDiv, driver.KeyKpDel,
			driver.KeyKpEnter, driver.KeyNums,
		},
	},
	{
		"Big Keys",
		[]string{
			driver.KeyF1, driver.KeyF2, driver.KeyF3, driver.KeyF4,
			driver.KeyF5, driver.KeyF6, driver.KeyF7, driver.KeyF8,
			driver.KeyF9, driver.KeyF10, driver.KeyF11, driver.KeyF12,
			"ESC", "TAB", "CAPS", "ENTER", "BACK", "SPACE",
			driver.KeyPrint, driver.KeyScrlk, driver.KeyPause,
			driver.KeyIns, driver.KeyHome, driver.KeyPgUp, driver.KeyDel, driver.KeyEnd, driver.KeyPgDn,
			driver.KeyArrUp, driver.KeyArrDown, driver.KeyArrLeft, driver.KeyArrRight,
			driver.KeyApp, driver.KeyK45, driver.KeyK56,
		},
	},
	{
		"Modifiers",
		[]string{
			driver.KeyCtrlL, driver.KeyShiftL, driver.KeyAltL, driver.KeyWinL,
			driver.KeyCtrlR, driver.KeyShiftR, driver.KeyAltR,
			driver.KeyFn1, driver.KeyFn2,
		},
	},
	{
		"Multimedia",
		[]string{
			driver.KeyVolUp, driver.KeyVolDn, driver.KeyMute,
			driver.KeyPlay, driver.KeyStop, driver.KeyPrev, driver.KeyNext,
		},
	},
	{
		"Mouse",
		[]string{
			driver.KeyMsL, driver.KeyMsR, driver.KeyMsM,
			driver.KeyMsScrU, driver.KeyMsScrD, driver.KeyMsScrL, driver.KeyMsScrR,
		},
	},
}

func NewActionPicker(onSelect func(action string), onNone func(), onClear func(), onCancel func()) *ActionPicker {
	var result *ActionPicker

	noneBtn := fynetool.NewButton("None", func() {
		if result != nil && result.Destroy != nil {
			result.Destroy()
		}
		if onNone != nil {
			onNone()
		}
	})
	clearBtn := fynetool.NewButton("Clear", func() {
		if result != nil && result.Destroy != nil {
			result.Destroy()
		}
		if onClear != nil {
			onClear()
		}
	})

	var items []fyne.CanvasObject
	items = append(items, container.NewGridWrap(fyne.NewSize(110, 52), noneBtn, clearBtn))

	for _, g := range actionGroups {
		group := g
		title := fynetool.NewLabel(group.Title)
		title.TextStyle = fyne.TextStyle{Bold: true}

		var row fyne.CanvasObject
		var btns []fyne.CanvasObject
		for _, a := range group.Actions {
			act := a
			btn := fynetool.NewButton(driver.ActionDisplayName(act), func() {
				if result != nil && result.Destroy != nil {
					result.Destroy()
				}
				if onSelect != nil {
					onSelect(act)
				}
			})
			btns = append(btns, btn)
		}
		row = container.NewGridWrap(fyne.NewSize(110, 52), btns...)
		items = append(items, title, row)
	}

	items = append(items, fynetool.NewLabel(""))
	scroll := container.NewVScroll(container.NewVBox(items...))

	cancelBtn := fynetool.NewButton("Cancel", func() {
		if result != nil && result.Destroy != nil {
			result.Destroy()
		}
		if onCancel != nil {
			onCancel()
		}
	})
	bottomRow := container.NewHBox(layout.NewSpacer(), cancelBtn)

	result = &ActionPicker{
		Content: container.NewBorder(nil, bottomRow, nil, nil, scroll),
		scroll:  scroll,
		Destroy: func() {},
	}
	return result
}
