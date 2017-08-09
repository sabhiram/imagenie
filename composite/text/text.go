package text

////////////////////////////////////////////////////////////////////////////////

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

////////////////////////////////////////////////////////////////////////////////

const (
	spacing = float64(1.1)
)

////////////////////////////////////////////////////////////////////////////////

var (
	f *truetype.Font
)

func SetupFont(fontpath string) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontpath)
	if err != nil {
		panic("Unable to read font file " + fontpath)
	}

	f, err = freetype.ParseFont(fontBytes)
	if err != nil {
		panic("Unable to parse font file " + fontpath)
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type Overlay struct {
	rotation   int
	xoff, yoff int
	size       float64
	value      string
	dpi        float64
	fontPath   string
	fg         color.Color
	bg         color.Color
}

func NewOverlay(ro, x, y, size, dpi int, fp string, fg, bg color.Color, value string) *Overlay {
	return &Overlay{
		rotation: ro,
		xoff:     x,
		yoff:     y,
		size:     float64(size),
		dpi:      float64(dpi),
		fontPath: fp,
		value:    value,
		fg:       fg,
		bg:       bg,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (o *Overlay) Render() (image.Image, int, int, int, error) {
	SetupFont(o.fontPath)

	scale := float64(o.dpi) / 72.0

	// Create an image to render the text on.
	img := image.NewRGBA(image.Rect(0, 0, int(o.size*spacing*float64(len(o.value))*scale), int(o.size*1.5*scale)))

	fg := image.NewUniform(o.fg)
	bg := image.NewUniform(o.bg)

	// Draw onto the image and setup the context.
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(o.dpi)
	c.SetFont(f)
	c.SetFontSize(o.size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	// Render the text on the context.
	ptl := freetype.Pt(0, int(c.PointToFixed(o.size)>>6))
	ptr, err := c.DrawString(o.value, ptl)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	xmax := ptr.X.Ceil() + 2
	ymax := ptr.Y.Ceil() + 2
	imgout := image.NewRGBA(image.Rect(0, 0, xmax, ymax))
	for i := 0; i < xmax; i++ {
		for j := 0; j < ymax; j++ {
			imgout.Set(i, j, img.At(i, j))
		}
	}
	return imgout, o.rotation, o.xoff, o.yoff, nil
}

////////////////////////////////////////////////////////////////////////////////
