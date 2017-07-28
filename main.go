////////////////////////////////////////////////////////////////////////////////

package main

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"path"
	"text/template"

	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/sabhiram/imagenie/composite"
	"github.com/sabhiram/imagenie/composite/image"
	"github.com/sabhiram/imagenie/composite/qr"
	"github.com/sabhiram/imagenie/composite/text"
)

////////////////////////////////////////////////////////////////////////////////

var (
	CLI = struct {
		outDir string   // output directory
		inFile string   // input file with pub and private keys
		args   []string // other args
	}{}
)

////////////////////////////////////////////////////////////////////////////////

// Function map definition, for any text template assistance.
var (
	funcMap = template.FuncMap{
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"add": func(a, b float64) float64 {
			return a + b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"div": func(a, b float64) float64 {
			return a / b
		},
		"firstHalf": func(s string) string {
			return s[:len(s)/2]
		},
		"secondHalf": func(s string) string {
			return s[len(s)/2:]
		},
		"precise8": func(a float64) string {
			return fmt.Sprintf("%.8f", a)
		},
		"precise4": func(a float64) string {
			return fmt.Sprintf("%.4f", a)
		},
	}
)

////////////////////////////////////////////////////////////////////////////////

type Overlay struct {
	Type     string `yaml:"type"`
	XOffset  int    `yaml:"xoffset"`
	YOffset  int    `yaml:"yoffset"`
	Size     int    `yaml:"size"`
	Template string `yaml:"template"`
	FgColor  string `yaml:"foreground"`
	BgColor  string `yaml:"background"`
}

// Hex parses a "html" hex color-string, either in the 3 "#f0c" or 6 "#ff1034" digits form.
// NOTE: This code has been borrowed and adapted from:
// 		 https://github.com/lucasb-eyer/go-colorful/blob/master/colors.go
func Hex(scol string) (color.Color, error) {
	format := "#%02x%02x%02x"
	factor := 1.0
	if len(scol) == 4 {
		format = "#%1x%1x%1x"
		factor = 16.0
	}

	var r, g, b uint8
	n, err := fmt.Sscanf(scol, format, &r, &g, &b)
	if err != nil {
		return color.RGBA{}, err
	}
	if n != 3 {
		return color.RGBA{}, fmt.Errorf("color: %v is not a hex-color", scol)
	}
	return color.RGBA{
		uint8(float64(r) * factor),
		uint8(float64(g) * factor),
		uint8(float64(b) * factor),
		255,
	}, nil
}

func getColor(c string, defaultColor color.Color) color.Color {
	switch c {
	case "black":
		return color.Black
	case "transparent":
		return color.Transparent
	case "white":
		return color.White
	case "":
		return defaultColor
	default:
		if c[0] == '#' {
			col, err := Hex(c)
			if err != nil {
				return defaultColor
			}
			return col
		}
	}
	return defaultColor
}

func (o *Overlay) GetRenderable(ctxt map[string]interface{}) (composite.Renderable, error) {
	var buf bytes.Buffer
	t := template.Must(template.New("output").Funcs(funcMap).Parse(o.Template))
	if err := t.Execute(&buf, ctxt); err != nil {
		log.Fatalf("Unable to execute template, error: %s\n", err.Error())
	}

	// Setup foreground color.
	fg := getColor(o.FgColor, color.Black)
	bg := getColor(o.BgColor, color.Transparent)

	switch o.Type {
	case "qr":
		return qr.NewOverlay(o.XOffset, o.YOffset, o.Size, buf.String(), fg, bg), nil
	case "text":
		return text.NewOverlay(o.XOffset, o.YOffset, o.Size, buf.String(), fg, bg), nil
	case "image":
		return image.NewOverlay(o.XOffset, o.YOffset, buf.String()), nil
	}
	return nil, fmt.Errorf("invalid renderable for overlay type: %s", o.Type)
}

////////////////////////////////////////////////////////////////////////////////

type Output struct {
	Prefix     string     `yaml:"prefix"`
	Background string     `yaml:"background"`
	Overlays   []*Overlay `yaml:"overlays"`
}

////////////////////////////////////////////////////////////////////////////////

type Specification struct {
	FontPath string                   `yaml:"fontpath"`
	Context  map[string]interface{}   `yaml:"context"`
	Items    []map[string]interface{} `yaml:"items"`
	Outputs  []*Output                `yaml:"outputs"`
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// TODO: If outdir does not exist, create it.

	CLI.args = flag.Args()
	if len(CLI.inFile) == 0 {
		log.Fatalf("specify input file with --infile!\n")
	}

	// Load the config file.
	raw, err := ioutil.ReadFile(CLI.inFile)
	if err != nil {
		log.Fatal(err)
	}

	var s Specification
	err = yaml.Unmarshal(raw, &s)
	if err != nil {
		log.Fatal(err)
	}

	// Load the global font that is specified in the sample file.
	text.SetupFont(s.FontPath)

	for index, m := range s.Items {
		// Build the context for each metadata item.
		ctxt := s.Context
		for k, v := range m {
			ctxt[k] = v
		}

		// Iterate all the outputs needed for this item.
		for _, output := range s.Outputs {
			// Build the set of renderables to build the ouput image.
			renderables := []composite.Renderable{}
			for _, overlay := range output.Overlays {
				renderable, err := overlay.GetRenderable(ctxt)
				if err != nil {
					log.Fatalf("Unable to get renderable for overlay. Error: %s\n", err.Error())
				}
				renderables = append(renderables, renderable)
			}

			// Generate the output image data.
			img, err := composite.BuildImage(output.Background, renderables)
			if err != nil {
				log.Fatalf("Unable to build image, error: %s\n", err.Error())
			}

			outFp := path.Join(CLI.outDir, fmt.Sprintf("%04d_%s.png", index, output.Prefix))
			outFd, err := os.OpenFile(outFp, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
			if err != nil {
				log.Fatalf("Unable to open file for write. %s\n", err.Error())
			}
			defer outFd.Close()

			if err := png.Encode(outFd, img); err != nil {
				log.Fatalf("Unable to encode file to png. %s\n", err.Error())
			}
		}
	}

	////////////////////////////////////////////////////////////

	// Build the inside image based on the background.
	// img, err := composite.BuildImage("./assets/PPIV_Inside.png", []composite.Renderable{
	// 	// Current value in PIVs and dollars.
	// 	text.NewOverlay(284, 250, 40, "Cash Value: $5.57"),
	// 	text.NewOverlay(284, 295, 40, "PIVS: 2.000000000"),

	// 	// QR code version of the private key.
	// 	qr.NewOverlay(456, 560, 180, "87eTMPUxKyiZTxEi6sdtUhRLcstWqCjva1ibxVdLg24zsH6o2XZ"),

	// 	// Split-up text version of private key.
	// 	text.NewOverlay(354, 560+182, 22, firstHalf("87eTMPUxKyiZTxEi6sdtUhRLcstWqCjva1ibxVdLg24zsH6o2XZ")),
	// 	text.NewOverlay(354, 560+182+20, 22, secondHalf("87eTMPUxKyiZTxEi6sdtUhRLcstWqCjva1ibxVdLg24zsH6o2XZ")),
	// })
	// if err != nil {
	// 	log.Fatalf("Error: Unable to build image: %s\n", err.Error())
	// }

}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetPrefix("")
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	flag.StringVar(&CLI.outDir, "outdir", "output", "output directory to put images")
	flag.StringVar(&CLI.inFile, "infile", "", "path to file that specifies keys to print")
	flag.Parse()
}

////////////////////////////////////////////////////////////////////////////////
