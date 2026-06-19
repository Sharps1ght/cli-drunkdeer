package widget

import (
	"image"
	"image/color"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var monoFontData []byte
var monoFontOnce sync.Once

func MonospaceFontData() []byte {
	monoFontOnce.Do(func() {
		path := findSystemMonospace()
		if path == "" {
			return
		}
		monoFontData, _ = os.ReadFile(path)
	})
	return monoFontData
}

func loadMonospaceFont() font.Face {
	data := MonospaceFontData()
	if len(data) == 0 {
		log.Println("No monospace TTF found, falling back to bitmap font")
		return basicfont.Face7x13
	}
	fnt, err := opentype.Parse(data)
	if err != nil {
		log.Printf("Failed to parse font: %v", err)
		return basicfont.Face7x13
	}
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    15,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Printf("Failed to create face: %v", err)
		return basicfont.Face7x13
	}
	return face
}

func findSystemMonospace() string {
	for _, q := range []string{"monospace:style=Bold", "monospace:bold", "monospace"} {
		out, err := exec.Command("fc-match", "--format=%{file}", q).Output()
		if err == nil {
			p := strings.TrimSpace(string(out))
			if p != "" {
				if _, err := os.Stat(p); err == nil {
					return p
				}
			}
		}
	}
	fallback := []string{
		"/usr/share/fonts/noto/NotoSansMono-Bold.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansMono-Bold.ttf",
		"/usr/share/fonts/noto/NotoSansMono-Regular.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansMono-Regular.ttf",
		"/usr/share/fonts/dejavu/DejaVuSansMono-Bold.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono-Bold.ttf",
		"/usr/share/fonts/dejavu/DejaVuSansMono.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		"/usr/share/fonts/TTF/DejaVuSansMono.ttf",
		"/usr/share/fonts/liberation/LiberationMono-Bold.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationMono-Bold.ttf",
		"/usr/share/fonts/liberation/LiberationMono-Regular.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf",
		"/usr/share/fonts/TTF/LiberationMono-Regular.ttf",
	}
	for _, p := range fallback {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

type KeyboardWidget struct {
	widget.BaseWidget
	layout             LayoutDef
	model              string
	selected           map[int]bool
	hoveredKey         int
	dragAnchor         fyne.Position
	dragging           bool
	dragStartSelected  map[int]bool
	raster             *canvas.Raster
	keyFace            font.Face
	defaultKeyColor    color.Color
	customColors       map[int]color.Color
	onSelectionChanged func()
}

func NewKeyboardWidget(model string) *KeyboardWidget {
	k := &KeyboardWidget{
		layout:         GetLayoutDef(model),
		model:          model,
		selected:       make(map[int]bool),
		hoveredKey:     -1,
		keyFace:        loadMonospaceFont(),
		defaultKeyColor: color.RGBA{0x2d, 0x2d, 0x2d, 0xff},
		customColors:   make(map[int]color.Color),
	}
	k.ExtendBaseWidget(k)
	return k
}

func parseHexColor(hex string) color.Color {
	if !strings.HasPrefix(hex, "#") || len(hex) != 7 {
		return color.RGBA{0x2d, 0x2d, 0x2d, 0xff}
	}
	r, _ := strconv.ParseUint(hex[1:3], 16, 8)
	g, _ := strconv.ParseUint(hex[3:5], 16, 8)
	b, _ := strconv.ParseUint(hex[5:7], 16, 8)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
}

func (k *KeyboardWidget) SetColors(defaultFill string, perKey map[int]string) {
	if defaultFill != "" {
		k.defaultKeyColor = parseHexColor(defaultFill)
	} else {
		k.defaultKeyColor = color.RGBA{0x2d, 0x2d, 0x2d, 0xff}
	}
	k.customColors = make(map[int]color.Color)
	for idx, hex := range perKey {
		k.customColors[idx] = parseHexColor(hex)
	}
	k.Refresh()
}

func (k *KeyboardWidget) SetIndividualColor(idx int, hex string) {
	k.customColors[idx] = parseHexColor(hex)
	k.Refresh()
}

func (k *KeyboardWidget) SelectedKeys() []int {
	var keys []int
	for v := range k.selected {
		keys = append(keys, v)
	}
	return keys
}

func (k *KeyboardWidget) ClearSelection() {
	k.selected = make(map[int]bool)
	if k.onSelectionChanged != nil {
		k.onSelectionChanged()
	}
	k.Refresh()
}

func (k *KeyboardWidget) SelectKeys(keys []int) {
	for _, v := range keys {
		k.selected[v] = true
	}
	if k.onSelectionChanged != nil {
		k.onSelectionChanged()
	}
	k.Refresh()
}

func (k *KeyboardWidget) SetOnSelectionChanged(f func()) {
	k.onSelectionChanged = f
}

func (k *KeyboardWidget) keyAt(x, y int) int {
	ld := k.layout
	yy := 0
	for _, row := range ld.Rows {
		xx := 0
		for _, key := range row {
			w := int(key.Width * ld.UnitSize)
			h := int(ld.RowHeight)
			if x >= xx && x < xx+w && y >= yy && y < yy+h {
				return key.Value
			}
			xx += w + int(ld.Gap)
		}
		yy += int(ld.RowHeight) + int(ld.Gap)
	}
	return -1
}

func (k *KeyboardWidget) Tapped(ev *fyne.PointEvent) {
	v := k.keyAt(int(ev.Position.X), int(ev.Position.Y))
	if v == -1 {
		return
	}
	if k.selected[v] {
		delete(k.selected, v)
	} else {
		k.selected[v] = true
	}
	if k.onSelectionChanged != nil {
		k.onSelectionChanged()
	}
	k.Refresh()
}

func (k *KeyboardWidget) MouseIn(ev *desktop.MouseEvent) {
	k.MouseMoved(ev)
}

func (k *KeyboardWidget) MouseMoved(ev *desktop.MouseEvent) {
	v := k.keyAt(int(ev.Position.X), int(ev.Position.Y))
	if v != k.hoveredKey {
		k.hoveredKey = v
		k.Refresh()
	}
}

func (k *KeyboardWidget) MouseOut() {
	if k.hoveredKey != -1 {
		k.hoveredKey = -1
		k.Refresh()
	}
}

func (k *KeyboardWidget) Dragged(ev *fyne.DragEvent) {
	if !k.dragging {
		k.dragging = true
		k.dragAnchor = ev.Position
		k.dragStartSelected = make(map[int]bool)
		for v := range k.selected {
			k.dragStartSelected[v] = true
		}
	}

	x1 := int(k.dragAnchor.X)
	y1 := int(k.dragAnchor.Y)
	x2 := int(ev.Position.X)
	y2 := int(ev.Position.Y)
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	k.selected = make(map[int]bool)
	for v := range k.dragStartSelected {
		k.selected[v] = true
	}
	ld := k.layout
	yy := 0
	for _, row := range ld.Rows {
		xx := 0
		for _, key := range row {
			kw := int(key.Width * ld.UnitSize)
			kh := int(ld.RowHeight)
			if rectsOverlap(x1, y1, x2-x1, y2-y1, xx, yy, kw, kh) {
				if k.dragStartSelected[key.Value] {
					delete(k.selected, key.Value)
				} else {
					k.selected[key.Value] = true
				}
			}
			xx += kw + int(ld.Gap)
		}
		yy += int(ld.RowHeight) + int(ld.Gap)
	}

	k.Refresh()
}

func (k *KeyboardWidget) DragEnd() {
	k.dragging = false
	k.dragStartSelected = nil
	if k.onSelectionChanged != nil {
		k.onSelectionChanged()
	}
	k.Refresh()
}

func rectsOverlap(ax, ay, aw, ah, bx, by, bw, bh int) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}

func (k *KeyboardWidget) CreateRenderer() fyne.WidgetRenderer {
	r := canvas.NewRaster(k.draw)
	k.raster = r
	return &keyboardRenderer{raster: r, widget: k}
}

func (k *KeyboardWidget) draw(w, h int) image.Image {
	ld := k.layout
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	y := 0
	for _, row := range ld.Rows {
		x := 0
		for _, key := range row {
			kw := int(key.Width * ld.UnitSize)
			kh := int(ld.RowHeight)

			var keyCol color.Color
			var textCol color.Color
			sel := k.selected[key.Value]
			hov := k.hoveredKey == key.Value
			switch {
			case sel && hov:
				keyCol = color.RGBA{0xff, 0x85, 0x66, 0xff}
				textCol = color.RGBA{0xff, 0xff, 0xff, 0xff}
			case sel:
				keyCol = color.RGBA{0xff, 0x63, 0x4d, 0xff}
				textCol = color.RGBA{0xff, 0xff, 0xff, 0xff}
			case hov:
				keyCol = color.RGBA{0x40, 0x40, 0x40, 0xff}
				textCol = color.RGBA{0xff, 0xff, 0xff, 0xff}
			default:
				if c, ok := k.customColors[key.Value]; ok {
					keyCol = c
				} else {
					keyCol = k.defaultKeyColor
				}
				textCol = color.RGBA{0xcc, 0xcc, 0xcc, 0xff}
			}

			drawRoundedRect(img, x, y, kw, kh, 4, keyCol)

			drawCenteredText(img, key.Name, x, y, kw, kh, textCol, k.keyFace)

			x += kw + int(ld.Gap)
		}
		y += int(ld.RowHeight) + int(ld.Gap)
	}

	return img
}

func drawRoundedRect(img *image.RGBA, x, y, w, h, r int, col color.Color) {
	rr, gg, bb, aa := col.RGBA()
	fill := color.RGBA{uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8), uint8(aa >> 8)}
	maxX := min(x+w, img.Bounds().Max.X)
	maxY := min(y+h, img.Bounds().Max.Y)
	r2 := r * r
	w1 := w - r - 1
	h1 := h - r - 1
	for py := y; py < maxY; py++ {
		for px := x; px < maxX; px++ {
			dx := px - x
			dy := py - y
			if dx < r && dy < r {
				if dx*dx+dy*dy > r2 {
					continue
				}
			} else if dx >= w-r && dy < r {
				dx2 := dx - w1
				if dx2*dx2+dy*dy > r2 {
					continue
				}
			} else if dx < r && dy >= h-r {
				dy2 := dy - h1
				if dx*dx+dy2*dy2 > r2 {
					continue
				}
			} else if dx >= w-r && dy >= h-r {
				dx2 := dx - w1
				dy2 := dy - h1
				if dx2*dx2+dy2*dy2 > r2 {
					continue
				}
			}
			img.Set(px, py, fill)
		}
	}
}

func drawCenteredText(img *image.RGBA, text string, x, y, w, h int, col color.Color, face font.Face) {
	m := face.Metrics()
	adv := font.MeasureString(face, text).Ceil()
	tx := x + (w-adv)/2
	ty := y + (h-m.Height.Ceil())/2 + m.Ascent.Ceil()

	px := fixed.I(tx)
	py := fixed.I(ty)
	faceColor := color.RGBA{0xff, 0xff, 0xff, 0xff}
	if c, ok := col.(color.RGBA); ok {
		faceColor = c
	}

	d := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(faceColor),
		Face: face,
		Dot:  fixed.Point26_6{X: px, Y: py},
	}
	d.DrawString(text)
}

type keyboardRenderer struct {
	raster *canvas.Raster
	widget *KeyboardWidget
}

func (r *keyboardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.raster}
}

func (r *keyboardRenderer) Layout(s fyne.Size) {
	r.raster.Resize(s)
}

func (r *keyboardRenderer) MinSize() fyne.Size {
	ld := r.widget.layout
	w := 0
	for _, row := range ld.Rows {
		rw := 0
		for _, k := range row {
			rw += int(k.Width*ld.UnitSize) + int(ld.Gap)
		}
		if rw > w {
			w = rw
		}
	}
	h := len(ld.Rows)*(int(ld.RowHeight)+int(ld.Gap)) - int(ld.Gap)
	return fyne.NewSize(float32(w), float32(h))
}

func (r *keyboardRenderer) Refresh() {
	r.raster.Refresh()
}

func (r *keyboardRenderer) ApplyTheme() {}
func (r *keyboardRenderer) BackgroundColor() color.Color {
	return color.Transparent
}
func (r *keyboardRenderer) Destroy() {}
