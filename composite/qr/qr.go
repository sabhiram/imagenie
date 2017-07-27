package qr

////////////////////////////////////////////////////////////////////////////////

import (
	"image"
	"image/color"

	qrcode "github.com/skip2/go-qrcode"
)

////////////////////////////////////////////////////////////////////////////////

type Overlay struct {
	xoff, yoff int
	width      int
	value      string
}

func NewOverlay(x, y, w int, value string) *Overlay {
	return &Overlay{
		xoff:  x,
		yoff:  y,
		width: w,
		value: value,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (o *Overlay) Render() (image.Image, int, int, error) {
	qr, err := qrcode.New(o.value, qrcode.Highest)
	if err != nil {
		return nil, 0, 0, err
	}

	qr.ForegroundColor = color.RGBA{0, 0, 0, 255}
	qr.BackgroundColor = color.RGBA{255, 255, 255, 0}

	// Note: We have a custom version of the `qr` library with one small
	// additional feature to scale the image up to avoid the QR having a large
	// offset. The scaling algorithm uses Lanczos3 as the filter, and can
	// be disabled if the results are undesirable.
	return qr.ImageNoPadding(o.width), o.xoff, o.yoff, nil
}

////////////////////////////////////////////////////////////////////////////////
