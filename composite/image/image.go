package image

////////////////////////////////////////////////////////////////////////////////

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

////////////////////////////////////////////////////////////////////////////////

const ()

////////////////////////////////////////////////////////////////////////////////

type Overlay struct {
	rotation   int
	xoff, yoff int
	value      string
}

func NewOverlay(ro, x, y int, value string) *Overlay {
	return &Overlay{
		rotation: ro,
		xoff:     x,
		yoff:     y,
		value:    value,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (o *Overlay) Render() (image.Image, int, int, int, error) {
	imgFd, err := os.Open(o.value)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	defer imgFd.Close()

	img, _, err := image.Decode(imgFd)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return img, o.rotation, o.xoff, o.yoff, nil
}

////////////////////////////////////////////////////////////////////////////////
