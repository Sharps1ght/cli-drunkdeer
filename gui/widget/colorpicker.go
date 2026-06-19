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

func makeSizedRaster(draw func(w, h int) image.Image, minW, minH int) *sizedRasterWidget {
	r := canvas.NewRaster(draw)
	s := &sizedRasterWidget{
		raster:  r,
		minSize: fyne.NewSize(float32(minW), float32(minH)),
	}
	s.ExtendBaseWidget(s)
	return s
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

	svRaster := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			sv := 1.0 - float64(y)/float64(h-1)
			for x := 0; x < w; x++ {
				ss := float64(x) / float64(w-1)
				cr, cg, cb := hsvToRGB(curHue, ss, sv)
				img.Set(x, y, color.RGBA{uint8(cr * 255), uint8(cg * 255), uint8(cb * 255), 0xff})
			}
		}
		cx := int(curSat * float64(w-1))
		cy := int((1.0 - curVal) * float64(h-1))
		for dx := -4; dx <= 4; dx++ {
			for dy := -4; dy <= 4; dy++ {
				if dx*dx+dy*dy > 4 {
					continue
				}
				px, py := cx+dx, cy+dy
				if px >= 0 && px < w && py >= 0 && py < h {
					img.Set(px, py, color.RGBA{0xff, 0xff, 0xff, 0xff})
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

	previewCombined := makeSizedRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		half := w / 2
		ac := hexToColor(initialHex)
		sc := hexToColor(curHex)
		for y := 0; y < h; y++ {
			for x := 0; x < half; x++ {
				img.Set(x, y, ac)
			}
			for x := half; x < w; x++ {
				img.Set(x, y, sc)
			}
		}
		return img
	}, 128, 32)

	var update func()
	var updating bool

	setSV := func(sv, vl float64) {
		curSat = math.Max(0, math.Min(1, sv))
		curVal = math.Max(0, math.Min(1, vl))
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

	svArea := &hBarTappable{
		raster:  svRaster,
		minSize: fyne.NewSize(130, 32),
		onDrag: func(x, y float64) {
			setSV(x, 1-y)
		},
	}
	svArea.ExtendBaseWidget(svArea)

	hueArea := &hBarTappable{
		raster:  hueRaster,
		minSize: fyne.NewSize(24, 32),
		onDrag: func(x, y float64) {
			setHue(x * 360)
		},
	}
	hueArea.ExtendBaseWidget(hueArea)

	hexEntry := fynetool.NewEntry()
	hexEntry.SetText(initialHex)
	hexEntry.PlaceHolder = "#FF00AA"
	hexWrapped := &minSizeWrap{inner: hexEntry, minsize: fyne.NewSize(90, 32)}
	hexWrapped.ExtendBaseWidget(hexWrapped)
	hexEntry.OnChanged = func(s string) {
		if updating {
			return
		}
		s = strings.TrimPrefix(s, "#")
		if len(s) == 6 {
			curHex = "#" + strings.ToUpper(s)
			curHue, curSat, curVal = hexToHSV(curHex)
			svRaster.Refresh()
			hueRaster.Refresh()
			previewCombined.Refresh()
			if onPreview != nil {
				onPreview(curHex)
			}
		}
	}

	update = func() {
		svRaster.Refresh()
		hueRaster.Refresh()
		previewCombined.Refresh()
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

	rightGroup := container.NewHBox(
		previewCombined,
		hexWrapped,
		selectBtn,
		cancelBtn,
	)

	barGrid := container.New(layout.NewGridLayout(2), svArea, hueArea)

	content := container.NewBorder(nil, nil,
		nil,
		rightGroup,
		barGrid,
	)

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
	return &simpleRender{obj: h.raster, minSize: h.minSize}
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

type minSizeWrap struct {
	fynetool.BaseWidget
	inner   fyne.CanvasObject
	minsize fyne.Size
}

func (m *minSizeWrap) CreateRenderer() fyne.WidgetRenderer {
	return &simpleRender{obj: m.inner, minSize: m.minsize}
}

type sizedRasterWidget struct {
	fynetool.BaseWidget
	raster  *canvas.Raster
	minSize fyne.Size
}

func (s *sizedRasterWidget) CreateRenderer() fyne.WidgetRenderer {
	return &simpleRender{obj: s.raster, minSize: s.minSize}
}

func (s *sizedRasterWidget) Refresh() {
	s.raster.Refresh()
}

type simpleRender struct {
	obj     fyne.CanvasObject
	minSize fyne.Size
}

func (r *simpleRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.obj}
}
func (r *simpleRender) Layout(s fyne.Size) { r.obj.Resize(s) }
func (r *simpleRender) MinSize() fyne.Size { return r.minSize }
func (r *simpleRender) Refresh()           { canvas.Refresh(r.obj) }
func (r *simpleRender) ApplyTheme()        {}
func (r *simpleRender) BackgroundColor() color.Color { return color.Transparent }
func (r *simpleRender) Destroy()           {}
