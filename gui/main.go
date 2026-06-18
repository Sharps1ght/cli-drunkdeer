package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetool "fyne.io/fyne/v2/widget"
	"github.com/2xxn/cli-drunkdeer/driver"
	kbwidget "github.com/2xxn/cli-drunkdeer/gui/widget"
	"github.com/sstallion/go-hid"
)

func findDrunkDeerDevice() *hid.DeviceInfo {
	products := []struct {
		VendorID  uint16
		ProductID uint16
		UsagePage uint16
		Usage     uint16
	}{
		{0x352D, 0x2383, 0xFF00, 0x00},
		{0x352D, 0x2382, 0xFF00, 0x00},
		{0x352D, 0x2384, 0xFF00, 0x00},
		{0x352D, 0x2386, 0xFF00, 0x00},
		{0x05AC, 0x024F, 0xFF00, 0x00},
	}

	var found *hid.DeviceInfo
	hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		for _, p := range products {
			if info.VendorID == p.VendorID && info.ProductID == p.ProductID &&
				info.UsagePage == p.UsagePage && info.Usage == p.Usage {
				cp := *info
				found = &cp
				return nil
			}
		}
		return nil
	})
	return found
}

func detectModel() string {
	devInfo := findDrunkDeerDevice()
	if devInfo == nil {
		return driver.KEYBOARD_A75
	}

	dev, err := hid.OpenPath(devInfo.Path)
	if err != nil {
		log.Printf("Failed to open device: %v", err)
		return driver.KEYBOARD_A75
	}
	defer dev.Close()

	controller := driver.NewDrunkDeerController(dev)
	identity := controller.GetIdentity()
	controller.Close()
	if identity != nil {
		return identity.KeyboardModel
	}
	return driver.KEYBOARD_A75
}

func main() {
	a := app.New()
	w := a.NewWindow("DrunkDeer Config")

	model := driver.KEYBOARD_G60
	log.Printf("Using keyboard model: %s", model)

	kb := kbwidget.NewKeyboardWidget(model)

	ok := fynetool.NewButton("OK", func() {
		fmt.Println("OK clicked")
	})
	cancel := fynetool.NewButton("Cancel", func() {
		fmt.Println("Cancel clicked")
	})
	buttons := container.NewCenter(container.NewHBox(ok, cancel))

	kbCentered := container.NewCenter(
		container.NewVBox(
			layout.NewSpacer(),
			kb,
			layout.NewSpacer(),
		),
	)

	content := container.NewBorder(nil, buttons, nil, nil, kbCentered)
	w.SetContent(content)

	pad := float32(80)
	butH := buttons.MinSize().Height
	w.Resize(fyne.NewSize(
		kb.MinSize().Width+pad,
		kb.MinSize().Height+pad+butH+float32(10),
	))
	w.CenterOnScreen()
	w.ShowAndRun()
}
