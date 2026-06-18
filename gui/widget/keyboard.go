package widget

import (
	"image"
	"image/color"
	"log"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func loadMonospaceFont() font.Face {
	path := findSystemMonospace()
	if path == "" {
		log.Println("No monospace TTF found, falling back to bitmap font")
		return basicfont.Face7x13
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read font %s: %v", path, err)
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
	layout     LayoutDef
	model      string
	selected   map[int]bool
	hoveredKey int
	raster     *canvas.Raster
	keyFace    font.Face
}

func NewKeyboardWidget(model string) *KeyboardWidget {
	k := &KeyboardWidget{
		layout:     GetLayoutDef(model),
		model:      model,
		selected:   make(map[int]bool),
		hoveredKey: -1,
		keyFace:    loadMonospaceFont(),
	}
	k.ExtendBaseWidget(k)
	return k
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
				keyCol = color.RGBA{0x2d, 0x2d, 0x2d, 0xff}
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
	for py := y; py < y+h && py < img.Bounds().Max.Y; py++ {
		for px := x; px < x+w && px < img.Bounds().Max.X; px++ {
			dx := px - x
			dy := py - y
			if dx < r && dy < r {
				if dx*dx+dy*dy > r*r {
					continue
				}
			}
			if dx >= w-r && dy < r {
				dx2 := dx - (w - r - 1)
				if dx2*dx2+dy*dy > r*r {
					continue
				}
			}
			if dx < r && dy >= h-r {
				dy2 := dy - (h - r - 1)
				if dx*dx+dy2*dy2 > r*r {
					continue
				}
			}
			if dx >= w-r && dy >= h-r {
				dx2 := dx - (w - r - 1)
				dy2 := dy - (h - r - 1)
				if dx2*dx2+dy2*dy2 > r*r {
					continue
				}
			}
			img.Set(px, py, color.RGBA{uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8), uint8(aa >> 8)})
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
	for _, k := range ld.Rows[0] {
		w += int(k.Width*ld.UnitSize) + int(ld.Gap)
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
