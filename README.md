# imagenie

[![Build Status](https://travis-ci.org/sabhiram/imagenie.svg?branch=master)](https://travis-ci.org/sabhiram/imagenie)

Batch overlay background images with text and QR codes.

## Why?

Often times there is a need to create a series of collateral where each paper, card, ticket etc, will require mostly static content with a few dynamic bits added to it.  These can be as simple as a name, phone number or a QR code that is consistently changing.  This tool aims to create a simple grammar by way of the input `yaml` file which will generate a batch of images with the specified overlays.

## Installation

```
go get github.com/sabhiram/imagenie
```

## Usage

Run the jobs specified in `sample.yaml` and generate images in the `outputs` directory.
```
imagenie -infile sample.yaml -outdir ./outputs
```

Please read the `./example/example.yaml` file on how to specify and configure jobs.

## Types of overlays

All overlays are required to be one of the following three types (which are shown in greater detail below):
1. `text`  - simple text based overlay
2. `image` - image overlay
3. `qr`    - qr code overlay

### Text

Similar to QR overlays, the text overlays will also require the X and Y offsets to position it.  However, the text overlays will use the specified global font path (.ttf file) to render the required text into a overlay.  Please note that the size supplied to the text overlay is an approximation of the font's size.

```yaml
      - type: text
        foreground: "white"
        background: "transparent"
        rotation: 90
        xoffset: 40
        yoffset: 40
        size: 40
        template: "Hi, I am {{ .gopher_name }}!"
```

### Image

An image overlay copies a target image at the specified offset into the background image.

```yaml
      - type: image
        xoffset: 320
        yoffset: 170
        template: ./assets/gopher.png
```

### QR

QR overlays are created by specifying an X and Y offset to the overlays expected location, and by setting a size in pixels for the QR code to span.  The value of the data fed to the QR code overlay generator will be converted to a `size` sized QR code with the higest data redundency.

```yaml
      - type: qr
        foreground: "black"
        background: "white"
        xoffset: 40
        yoffset: 140
        size: 256
        template: "Gopher {{ .gopher_name }} has ID : {{ .gopher_id }}"
```

## Overlay options

All "jobs" start off with a background image.  This is the base image which will be built upon.  All overlays have the following optional properties:
1. `foreground` - text / qr color, invalid for image overlays
2. `background` - background color, invalid for image overlays

The default foreground color is black, and the default background is transparent. You can specify colors for the `foreground` and `background` in the following ways:

1. `black`, `white` and `transparent` - are valid values.
2. Any hex value in the form of "#FFFFFF" (white)
3. Any hex value in the form of "#F00" (red)

You can additionally specify the rotation that needs to be applied to a given overlay.  All rotations will be applied before the offsetting of x and y, and the rotations will be counter-clockwise.  Valid values include any number from 0-360.  The default rotation will be 0 degrees.

## Sample Usage

For a detailed example, check out the `./example/README.md` file, as well as the accompanying `./example/example.yaml` file.

## Notes:

Cross compile for windows:

```
GOOS=windows GOARCH=386 go build -o imagenie.exe .
```