package main

import (
	"./encoder"
	"./generators/pixelated"
	"github.com/codegangsta/cli"
	"os"
	"fmt"
)

func main() {
	app := cli.NewApp()
	app.Name = "Avatarme"
	app.Usage = "Generates an unique avatar for the given string"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "output, o", Value: "output.png", Usage: "path of the output file"},
		cli.IntFlag{Name: "size, s", Value: 256, Usage: "side length of the generated image (in px). Will be ronded to a multiple of 8"},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) <= 0 {
			fmt.Println("Error: You need to supply a string to encode")
			return
		}
		text := c.Args()[0]
		img := pixelated.BuildImage(text, c.Int("size"))
		encoder.ExportImage(img, c.String("output"))
	}

	app.Run(os.Args)
}
