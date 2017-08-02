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

////////////////////////////////////////////////////////////////////////////////

func defaultIntValue(v, def int) int {
	if v == 0 {
		return def
	}
	return v
}

////////////////////////////////////////////////////////////////////////////////

// OverlayOpts specifies the possible options to configure a given overlay type.
// The types of overlays that each option applies to are specified in the
// comment to the right of the declaration.
type OverlayOpts struct {
	Type     string `yaml:"type"`       // Image, QR, Text
	XOffset  int    `yaml:"xoffset"`    // Image, QR, Text
	YOffset  int    `yaml:"yoffset"`    // Image, QR, Text
	Size     int    `yaml:"size"`       // Image, QR, Text
	Dpi      int    `yaml:"dpi"`        // Text
	Template string `yaml:"template"`   // Image, QR, Text
	FgColor  string `yaml:"foreground"` // QR, Text
	BgColor  string `yaml:"background"` // QR, Text
}

// GetRenderable returns a `Renderable` interface based on the underlying overlay
// options.
func (o *OverlayOpts) GetRenderable(ctxt map[string]interface{}) (composite.Renderable, error) {
	var buf bytes.Buffer
	t := template.Must(template.New("output").Funcs(funcMap).Parse(o.Template))
	if err := t.Execute(&buf, ctxt); err != nil {
		log.Fatalf("Unable to execute template, error: %s\n", err.Error())
	}

	// The templated value is the string to either print or QR in the case
	// of those overlay types.  In the case of the image type, it is a path to
	// the image to inject to allow for a dynamic range of images to be used.
	tv := buf.String()

	// Default values in case they are not configured
	xo := o.XOffset                   // Default: 0
	yo := o.YOffset                   // Default: 0
	sz := defaultIntValue(o.Size, 12) // Default: 12 "pt"
	dp := defaultIntValue(o.Dpi, 72)  // Default: 72 dpi
	fg := getColor(o.FgColor, color.Black)
	bg := getColor(o.BgColor, color.Transparent)

	switch o.Type {
	case "qr":
		return qr.NewOverlay(xo, yo, sz, fg, bg, tv), nil
	case "text":
		return text.NewOverlay(xo, yo, sz, dp, fg, bg, tv), nil
	case "image":
		return image.NewOverlay(xo, yo, tv), nil
	}
	return nil, fmt.Errorf("invalid renderable for overlay type: %s", o.Type)
}

////////////////////////////////////////////////////////////////////////////////

// Output represents a single job to be done for a given background image, and
// the list of overlays that are to be applied to the same.
type Output struct {
	Prefix     string         `yaml:"prefix"`
	Background string         `yaml:"background"`
	Overlays   []*OverlayOpts `yaml:"overlays"`
}

////////////////////////////////////////////////////////////////////////////////

// Config represents the config file needed to run the program.
type Config struct {
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

	var cfg Config
	err = yaml.Unmarshal(raw, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Load the global font that is specified in the sample file.
	text.SetupFont(cfg.FontPath)

	// Iterate through all the "jobs" that we need to carry out.
	for _, output := range cfg.Outputs {
		log.Printf("Processing job with prefix: %s (%s)\n", output.Prefix, output.Background)
		for index, m := range cfg.Items {
			log.Printf("  Processing item #%d\n", index+1)

			// Build the context for each metadata item.
			ctxt := cfg.Context
			for k, v := range m {
				ctxt[k] = v
			}

			// Build the set of renderables to build the ouput image.
			renderables := []composite.Renderable{}
			for idx, overlay := range output.Overlays {
				log.Printf("    * adding %s overlay at index %d\n", overlay.Type, idx+1)
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
			log.Printf("  --> Generated output file: %s\n", outFp)
		}
		log.Printf("\n")
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetPrefix("")
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	flag.StringVar(&CLI.outDir, "outdir", "output", "output directory to put images")
	flag.StringVar(&CLI.outDir, "o", "output", "output directory to put images (short)")
	flag.StringVar(&CLI.inFile, "infile", "", "path to file that specifies keys to print")
	flag.StringVar(&CLI.inFile, "i", "", "path to file that specifies keys to print (short)")
	flag.Parse()
}

////////////////////////////////////////////////////////////////////////////////
