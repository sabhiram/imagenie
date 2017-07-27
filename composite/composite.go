package composite

////////////////////////////////////////////////////////////////////////////////
/*

Composite implements a single method that accepts a background image (path), and
a list of renderables to render.

Each renderable implements a "Render" function as part of adhering to the
Renderable interface defined below.

*/
////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"image"
	"os"
)

////////////////////////////////////////////////////////////////////////////////

type Renderable interface {
	Render() (image.Image, int, int, error)
}

////////////////////////////////////////////////////////////////////////////////

func BuildImage(bgpath string, items []Renderable) (*image.RGBA, error) {
	baseImgFd, err := os.Open(bgpath)
	if err != nil {
		return nil, err
	}
	defer baseImgFd.Close()

	baseImg, _, err := image.Decode(baseImgFd)
	if err != nil {
		return nil, err
	}

	// Create an output image, copy each pixel from the background to the temp
	// image so that we can build up each layer of the overlays.
	bounds := baseImg.Bounds()
	out := image.NewRGBA(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			out.Set(x, y, baseImg.At(x, y))
		}
	}

	// Overlay each renderable on top of the image.
	for i, item := range items {
		fmt.Printf("Applying overlay index: %d for item: %#v\n", i, item)
		img, xoff, yoff, err := item.Render()
		if err != nil {
			return nil, err
		}

		inbounds := img.Bounds()
		w, h := inbounds.Max.X, inbounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if x+xoff < bounds.Max.X && y+yoff < bounds.Max.Y {
					_, _, _, a := img.At(x, y).RGBA()
					if a >= 3*0xffff/4 {
						out.Set(x+xoff, y+yoff, img.At(x, y))
					}
				}
			}
		}
	}

	return out, nil
}

////////////////////////////////////////////////////////////////////////////////
