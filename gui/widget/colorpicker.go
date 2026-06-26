package widget

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	fynetool "fyne.io/fyne/v2/widget"
)

func hsvToRGB(h, s, v float64) (float64, float64, float64) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return r + m, g + m, b + m
}

func hexToHSV(hex string) (h, s, v float64) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) < 6 {
		return 0, 0, 1
	}
	var ri, gi, bi uint8
	fmt.Sscanf(hex[:2], "%02x", &ri)
	fmt.Sscanf(hex[2:4], "%02x", &gi)
	fmt.Sscanf(hex[4:6], "%02x", &bi)
	rf := float64(ri) / 255
	gf := float64(gi) / 255
	bf := float64(bi) / 255
	maxV := math.Max(rf, math.Max(gf, bf))
	minV := math.Min(rf, math.Min(gf, bf))
	delta := maxV - minV
	v = maxV
	if maxV == 0 {
		s = 0
	} else {
		s = delta / maxV
	}
	if delta == 0 {
		h = 0
	} else if maxV == rf {
		h = 60 * math.Mod((gf-bf)/delta, 6)
	} else if maxV == gf {
		h = 60 * ((bf-rf)/delta + 2)
	} else {
		h = 60 * ((rf-gf)/delta + 4)
	}
	if h < 0 {
		h += 360
	}
	return
}

func hexToColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) < 6 {
		return color.RGBA{0, 0, 0, 0xff}
	}
	var r, g, b uint8
	fmt.Sscanf(hex[:2], "%02x", &r)
	fmt.Sscanf(hex[2:4], "%02x", &g)
	fmt.Sscanf(hex[4:6], "%02x", &b)
	return color.RGBA{r, g, b, 0xff}
}

type ColorPicker struct {
	Content fyne.CanvasObject
	Destroy func()
}

func NewColorPicker(
	initialHex string,
	onPreview func(hex string),
	onSelect func(hex string),
	onCancel func(),
) *ColorPicker {
	h, s, v := hexToHSV(initialHex)
	var curHex = initialHex
	var curHue = h
	var curSat = s
	var curVal = v

	satRaster := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			ss := float64(x) / float64(w-1)
			cr, cg, cb := hsvToRGB(curHue, ss, curVal)
			col := color.RGBA{uint8(cr * 255), uint8(cg * 255), uint8(cb * 255), 0xff}
			for y := 0; y < h; y++ {
				img.Set(x, y, col)
			}
		}
		cx := int(curSat * float64(w-1))
		for dx := -2; dx <= 2; dx++ {
			for y := 0; y < h; y++ {
				px := cx + dx
				if px >= 0 && px < w {
					img.Set(px, y, color.RGBA{0xff, 0xff, 0xff, 0xff})
				}
			}
		}
		return img
	})

	valRaster := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			sv := float64(x) / float64(w-1)
			cr, cg, cb := hsvToRGB(curHue, curSat, sv)
			col := color.RGBA{uint8(cr * 255), uint8(cg * 255), uint8(cb * 255), 0xff}
			for y := 0; y < h; y++ {
				img.Set(x, y, col)
			}
		}
		cx := int(curVal * float64(w-1))
		for dx := -2; dx <= 2; dx++ {
			for y := 0; y < h; y++ {
				px := cx + dx
				if px >= 0 && px < w {
					img.Set(px, y, color.RGBA{0xff, 0xff, 0xff, 0xff})
				}
			}
		}
		return img
	})

	hueRaster := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				hue := float64(x) / float64(w-1) * 360
				cr, cg, cb := hsvToRGB(hue, 1, 1)
				img.Set(x, y, color.RGBA{uint8(cr * 255), uint8(cg * 255), uint8(cb * 255), 0xff})
			}
		}
		cx := int(curHue / 360 * float64(w-1))
		for dx := -2; dx <= 2; dx++ {
			for dy := 0; dy < h; dy++ {
				px := cx + dx
				if px >= 0 && px < w {
					img.Set(px, dy, color.RGBA{0xff, 0xff, 0xff, 0xff})
				}
			}
		}
		return img
	})

	refColor := canvas.NewRectangle(hexToColor(initialHex))
	refColor.SetMinSize(fyne.NewSize(156, 38))

	updatePreview := func() {
	}

	var update func()
	var updating bool

	setSat := func(s float64) {
		curSat = math.Max(0, math.Min(1, s))
		cr, cg, cb := hsvToRGB(curHue, curSat, curVal)
		curHex = fmt.Sprintf("#%02X%02X%02X", uint8(cr*255), uint8(cg*255), uint8(cb*255))
		update()
	}

	setVal := func(v float64) {
		curVal = math.Max(0, math.Min(1, v))
		cr, cg, cb := hsvToRGB(curHue, curSat, curVal)
		curHex = fmt.Sprintf("#%02X%02X%02X", uint8(cr*255), uint8(cg*255), uint8(cb*255))
		update()
	}

	setHue := func(h float64) {
		curHue = math.Max(0, math.Min(360, h))
		cr, cg, cb := hsvToRGB(curHue, curSat, curVal)
		curHex = fmt.Sprintf("#%02X%02X%02X", uint8(cr*255), uint8(cg*255), uint8(cb*255))
		update()
	}

	barSize := fyne.NewSize(156, 38)

	satArea := &hBarTappable{
		raster:  satRaster,
		minSize: barSize,
		onDrag: func(x, y float64) {
			setSat(x)
		},
	}
	satArea.ExtendBaseWidget(satArea)

	valArea := &hBarTappable{
		raster:  valRaster,
		minSize: barSize,
		onDrag: func(x, y float64) {
			setVal(x)
		},
	}
	valArea.ExtendBaseWidget(valArea)

	hueArea := &hBarTappable{
		raster:  hueRaster,
		minSize: barSize,
		onDrag: func(x, y float64) {
			setHue(x * 360)
		},
	}
	hueArea.ExtendBaseWidget(hueArea)

	hexEntry := fynetool.NewEntry()
	hexEntry.SetText(initialHex)
	hexEntry.PlaceHolder = "#FF00AA"
	hexWrapped := NewMinSizeWrap(hexEntry, fyne.NewSize(108, 38))
	hexEntry.OnChanged = func(s string) {
		if updating {
			return
		}
		s = strings.TrimPrefix(s, "#")
		if len(s) == 6 {
			curHex = "#" + strings.ToUpper(s)
			curHue, curSat, curVal = hexToHSV(curHex)
			satRaster.Refresh()
			valRaster.Refresh()
			hueRaster.Refresh()
			updatePreview()
			if onPreview != nil {
				onPreview(curHex)
			}
		}
	}

	update = func() {
		satRaster.Refresh()
		valRaster.Refresh()
		hueRaster.Refresh()
		updatePreview()
		updating = true
		hexEntry.SetText(curHex)
		updating = false
		if onPreview != nil {
			onPreview(curHex)
		}
	}

	var result *ColorPicker

	selectBtn := fynetool.NewButton("Select", func() {
		if result != nil && result.Destroy != nil {
			result.Destroy()
		}
		if onSelect != nil {
			onSelect(curHex)
		}
	})
	cancelBtn := fynetool.NewButton("Cancel", func() {
		if result != nil && result.Destroy != nil {
			result.Destroy()
		}
		if onCancel != nil {
			onCancel()
		}
	})

	rightGroup := container.NewVBox(
		refColor,
		hexWrapped,
		container.NewHBox(selectBtn, cancelBtn),
		layout.NewSpacer(),
	)

	bars := container.NewVBox(hueArea, satArea, valArea, layout.NewSpacer())

	content := container.NewHBox(bars, rightGroup)

	result = &ColorPicker{
		Content: content,
		Destroy: func() {},
	}
	return result
}

type hBarTappable struct {
	fynetool.BaseWidget
	raster  *canvas.Raster
	onDrag  func(x, y float64)
	minSize fyne.Size
}

func (h *hBarTappable) CreateRenderer() fyne.WidgetRenderer {
	return &spriteRenderer{inner: h.raster, minSize: h.minSize}
}

func (h *hBarTappable) Tapped(ev *fyne.PointEvent) {
	size := h.raster.Size()
	if size.Width > 0 && size.Height > 0 {
		h.onDrag(float64(ev.Position.X)/float64(size.Width), float64(ev.Position.Y)/float64(size.Height))
	}
}

func (h *hBarTappable) Dragged(ev *fyne.DragEvent) {
	size := h.raster.Size()
	if size.Width > 0 && size.Height > 0 {
		h.onDrag(float64(ev.Position.X)/float64(size.Width), float64(ev.Position.Y)/float64(size.Height))
	}
}

func (h *hBarTappable) DragEnd() {}


