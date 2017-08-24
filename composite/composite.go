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
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/disintegration/imaging"
)

////////////////////////////////////////////////////////////////////////////////

type Renderable interface {
	Render() (image.Image, int, int, int, error)
}

////////////////////////////////////////////////////////////////////////////////

func BuildImage(bgpath, ofpath, offmt string, items []Renderable) error {
	baseImgFd, err := os.Open(bgpath)
	if err != nil {
		return err
	}
	defer baseImgFd.Close()

	baseImg, _, err := image.Decode(baseImgFd)
	if err != nil {
		return err
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
		primg, rot, xoff, yoff, err := item.Render()
		if err != nil {
			return err
		}

		img := primg
		if rot > 0 && rot < 360 {
			img = imaging.Rotate(primg, float64(rot), color.Transparent)
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

	// Emit the file as an image in the specified output format, location.
	outfd, err := os.OpenFile(ofpath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer outfd.Close()

	switch strings.ToLower(offmt) {
	case "png":
		if err := png.Encode(outfd, out); err != nil {
			return err
		}
	case "jpeg", "jpg":
		if err := jpeg.Encode(outfd, out, nil); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a valid output format type", offmt)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func BuildImageWithMagick(binspath, bgpath, ofpath, offmt, ofcs string, items []Renderable) error {
	// Create the output image as a copy of the background.
	cmd := fmt.Sprintf("%s ( +clone ) -composite %s", bgpath, ofpath)
	cmdCopy := exec.Command(path.Join(binspath, "convert"), strings.Split(cmd, " ")...)
	_, err := cmdCopy.CombinedOutput()
	if err != nil {
		return err
	}

	// Overlay each renderable on top of the image.
	for _, item := range items {
		primg, rot, xoff, yoff, err := item.Render()
		if err != nil {
			return err
		}
		img := primg
		if rot > 0 && rot < 360 {
			img = imaging.Rotate(primg, float64(rot), color.Transparent)
		}

		tempImgPath := "tmp.png"
		tempFd, err := os.OpenFile(tempImgPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer tempFd.Close()

		if err := png.Encode(tempFd, img); err != nil {
			return err
		}

		// Imagemagick only understands "rgb" as a valid colorspace
		if ofcs == "rgba" || ofcs == "rgb" {
			ofcs = "rgb"
		} else if ofcs == "cmyk" {
			ofcs = "cmyk"
		} else {
			return fmt.Errorf("%s is an invalid colorspace!", ofcs)
		}

		cmd = fmt.Sprintf("-colorspace %s -compose atop -geometry +%d+%d %s %s %s", ofcs, xoff, yoff, tempImgPath, ofpath, ofpath)
		cmd1 := exec.Command(path.Join(binspath, "composite"), strings.Split(cmd, " ")...)
		_, err = cmd1.CombinedOutput()
		if err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
