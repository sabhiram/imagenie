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
	"image"
	"image/color"
	"os"

	_ "image/jpeg"
	_ "image/png"
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
	for _, item := range items {
		img, xoff, yoff, err := item.Render()
		if err != nil {
			return nil, err
		}

		inbounds := img.Bounds()
		w, h := inbounds.Max.X, inbounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if x+xoff < bounds.Max.X && y+yoff < bounds.Max.Y {
					rf, gf, bf, af := img.At(x, y).RGBA()
					rb, gb, bb, ab := out.At(x+xoff, y+yoff).RGBA()
					alpha := float64(af) / float64(0xffff)
					beta := 1.0 - alpha

					// Some hacky alpha blending - revisit later :)
					c := color.RGBA{
						uint8((float64(rf)*alpha + float64(rb)*beta) / 256),
						uint8((float64(gf)*alpha + float64(gb)*beta) / 256),
						uint8((float64(bf)*alpha + float64(bb)*beta) / 256),
						uint8((float64(af)*alpha + float64(ab)*beta) / 256),
					}
					out.Set(x+xoff, y+yoff, c)
				}
			}
		}
	}

	return out, nil
}

////////////////////////////////////////////////////////////////////////////////
