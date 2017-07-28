# imagenie example

Ok, so the inter-go-lactic foundation has tasked you with providing a million images of a million gophers who have visited planet `plutgo`.  They require you to generate an image of the gopher, stamped with the gopher's custom QR code, and a bunch of other information like when the gopher visited the planet and other such intriguing facts.  This all sounds well and good, but knowing the foundation, you are certain that they will make you move a bunch of things left and right until they are "right enough" and ask you to change the wording for all million right after you finish the job.  Goodthing you are reading this example!

This directory contains all the assets needed to run the `imagenie` tool, and overlay this background:

![background](https://raw.githubusercontent.com/sabhiram/imagenie/master/example/assets/bg.jpeg)

With this sexy gopher:

![gopher](https://raw.githubusercontent.com/sabhiram/imagenie/master/example/assets/gopher.png)

... and a few other pieces of information, to generate this:

![gopher](https://raw.githubusercontent.com/sabhiram/imagenie/master/example/assets/output.png)

You can also generate a whole batch of them (like a million)!

## Running this is simple

From this directory:
```
$ go run ../main.go -infile example.yaml -outdir .
```

This will execute the example and generate a bunch of images in the current directory. Feel free to play with the (hopefully) well documented `example.yaml` file in this directory.

## Contents

### assets

This directory contains all required collateral for generating the above "job".  Things like the background image, the superimposed gopher, any font files live here.

### example.yaml

For detailed information on the grammar that we can use, and for a few simple examples on the templating that is available, check out this file.  This is essentially the core of the program and defines what gets overlayed over the background and also defines the list of inputs to permute over.
