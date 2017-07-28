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
	dpi     = float64(72)
	spacing = float64(1.1)
)

////////////////////////////////////////////////////////////////////////////////

var (
	f *truetype.Font
)

////////////////////////////////////////////////////////////////////////////////

type Overlay struct {
	xoff, yoff int
	size       float64
	value      string
	fg         color.Color
	bg         color.Color
}

func NewOverlay(x, y, size int, value string, fg, bg color.Color) *Overlay {
	return &Overlay{
		xoff:  x,
		yoff:  y,
		size:  float64(size),
		value: value,
		fg:    fg,
		bg:    bg,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (o *Overlay) Render() (image.Image, int, int, error) {
	// Create an image to render the text on.
	img := image.NewRGBA(image.Rect(0, 0, int(o.size*spacing*float64(len(o.value))), int(o.size*1.5)))

	fg := image.NewUniform(o.fg)
	bg := image.NewUniform(o.bg)

	// Draw onto the image and setup the context.
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(o.size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	// Render the text on the context.
	pt := freetype.Pt(45, int(c.PointToFixed(o.size)>>6))
	if _, err := c.DrawString(o.value, pt); err != nil {
		return nil, 0, 0, err
	}
	return img, o.xoff, o.yoff, nil
}

////////////////////////////////////////////////////////////////////////////////

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
