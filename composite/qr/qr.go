package qr

////////////////////////////////////////////////////////////////////////////////

import (
	"image"
	"image/color"

	qrcode "github.com/skip2/go-qrcode"
)

////////////////////////////////////////////////////////////////////////////////

type Overlay struct {
	rotation   int
	xoff, yoff int
	width      int
	value      string
	fg         color.Color
	bg         color.Color
}

func NewOverlay(ro, x, y, w int, fg, bg color.Color, value string) *Overlay {
	return &Overlay{
		rotation: ro,
		xoff:     x,
		yoff:     y,
		width:    w,
		value:    value,
		fg:       fg,
		bg:       bg,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (o *Overlay) Render() (image.Image, int, int, int, error) {
	qr, err := qrcode.New(o.value, qrcode.Highest)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	qr.ForegroundColor = o.fg
	qr.BackgroundColor = o.bg

	// Note: We have a custom version of the `qr` library with one small
	// additional feature to scale the image up to avoid the QR having a large
	// offset. The scaling algorithm uses Lanczos3 as the filter, and can
	// be disabled if the results are undesirable.
	return qr.ImageNoPadding(o.width), o.rotation, o.xoff, o.yoff, nil
}

////////////////////////////////////////////////////////////////////////////////
