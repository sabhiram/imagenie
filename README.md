# imagenie

Batch overlay background images with text and QR codes.

## Why?

Often times there is a need to create a series of collateral where each paper, card, ticket etc, will require mostly static content with a few dynamic bits added to it.  These can be as simple as a name, phone number or a QR code that is consistently changing.  This tool aims to create a simple grammer by way of the input `yaml` file which will generate a batch of images with the specified overlays.

## Types of overlays

All "jobs" start off with a background image.  This is the base image which will be built upon.  The only two types of overlays supported at the moment are:

### QR

QR overlays are created by specifying an X and Y offset to the overlays expected location, and by setting a size in pixels for the QR code to span.  The value of the data fed to the QR code overlay generator will be converted to a `size` sized QR code with the higest data redundency.

### Text

Similar to QR overlays, the text overlays will also require the X and Y offsets to position it.  However, the text overlays will use the specified global font path (.ttf file) to render the required text into a overlay.  Please note that the size supplied to the text overlay is an approximation of the font's size.

## Sample Usage



## Installation

TODO

## Notes:

Cross compile for windows:

```
GOOS=windows GOARCH=386 go build -o imagenie.exe .
```