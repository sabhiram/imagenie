////////////////////////////////////////////////////////////////////////////////

package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/sabhiram/PIVX-WalletGen/composite"
	"github.com/sabhiram/PIVX-WalletGen/composite/qr"
	"github.com/sabhiram/PIVX-WalletGen/composite/text"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

////////////////////////////////////////////////////////////////////////////////

type Overlay struct {
	Type     string `yaml:"type"`
	XOffset  int    `yaml:"xoffset"`
	YOffset  int    `yaml:"yoffset"`
	Size     int    `yaml:"size"`
	Template string `yaml:"template"`
}

type Output struct {
	Prefix     string     `yaml:"prefix"`
	Background string     `yaml:"background"`
	Overlays   []*Overlay `yaml:"overlays"`
}

type Specification struct {
	Context  map[string]interface{}   `yaml:"context"`
	MetaData []map[string]interface{} `yaml:"metadata"`
	Outputs  []*Output                `yaml:"outputs"`
}

func testYaml() {
	// Load the config file.
	raw, err := ioutil.ReadFile("sample.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var s Specification
	err = yaml.Unmarshal(raw, &s)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("GOT STRUCT: %#v\n", s)
}

////////////////////////////////////////////////////////////////////////////////

var (
	CLI = struct {
		outDir string   // output directory
		inFile string   // input file with pub and private keys
		args   []string // other args
	}{}
)

////////////////////////////////////////////////////////////////////////////////

type PPiv struct {
	publicKey  string
	privateKey string
	expiryDate string
	pivValue   float64
	usdValue   float64
}

func NewPPiv(pub, priv, expiry string, pivs, usd float64) *PPiv {
	return &PPiv{
		publicKey:  pub,
		privateKey: priv,
		expiryDate: expiry,
		pivValue:   pivs,
		usdValue:   usd,
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	CLI.args = flag.Args()
	if len(CLI.inFile) == 0 {
		log.Fatalf("specify input file with --infile!\n")
	}

	////////////////////////////////////////////////////////////

	// Build the inside image based on the background.
	img, err := composite.BuildImage("./assets/PPIV_Inside.png", []composite.Renderable{
		// Current value in PIVs and dollars.
		text.NewOverlay(284, 250, 40, "Cash Value: $5.57"),
		text.NewOverlay(284, 295, 40, "PIVS: 2.000000000"),

		// QR code version of the private key.
		qr.NewOverlay(456, 560, 180, "87eTMPUxKyiZTxEi6sdtUhRLcstWqCjva1ibxVdLg24zsH6o2XZ"),

		// Split-up text version of private key.
		text.NewOverlay(354, 560+182, 22, firstHalf("87eTMPUxKyiZTxEi6sdtUhRLcstWqCjva1ibxVdLg24zsH6o2XZ")),
		text.NewOverlay(354, 560+182+20, 22, secondHalf("87eTMPUxKyiZTxEi6sdtUhRLcstWqCjva1ibxVdLg24zsH6o2XZ")),
	})
	if err != nil {
		log.Fatalf("Error: Unable to build image: %s\n", err.Error())
	}

	////////////////////////////////////////////////////////////

	outFd, err := os.OpenFile("output.png", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Unable to open file for write. %s\n", err.Error())
	}
	defer outFd.Close()

	if err := png.Encode(outFd, img); err != nil {
		log.Fatalf("Unable to encode file to png. %s\n", err.Error())
	}

	testYaml()
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
