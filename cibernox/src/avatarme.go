package main

import (
	"crypto/md5"
	"fmt"
	"github.com/codegangsta/cli"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"regexp"
)

const (
	enlarger        string = "md5-enlarger"
	avatarSide      int    = 8
	signatureLength int    = avatarSide * avatarSide
)

var (
	extensionRegexp *regexp.Regexp = regexp.MustCompile(`\.(png|jpe?g)$`)
)

var palette map[byte]color.RGBA = map[byte]color.RGBA{
	48:  color.RGBA{255, 255, 255, 255}, // 0 => white,
	49:  color.RGBA{213, 0, 0, 255},     // 1 => red,
	50:  color.RGBA{255, 255, 255, 255}, // 2 => white,
	51:  color.RGBA{255, 255, 255, 255}, // 3 => white,
	52:  color.RGBA{255, 76, 0, 255},    // 4 => orange,
	53:  color.RGBA{255, 255, 255, 255}, // 5 => white,
	54:  color.RGBA{255, 255, 255, 255}, // 6 => white,
	55:  color.RGBA{255, 255, 11, 255},  // 7 => yellow,
	56:  color.RGBA{255, 76, 0, 255},    // 8 => orange,
	57:  color.RGBA{213, 0, 0, 255},     // 9 => red,
	97:  color.RGBA{239, 0, 113, 255},   // a => magenta,
	98:  color.RGBA{54, 0, 151, 255},    // b => purple,
	99:  color.RGBA{0, 0, 205, 255},     // c => blue,
	100: color.RGBA{0, 152, 232, 255},   // d => cyan,
	101: color.RGBA{26, 176, 0, 255},    // e => green,
	102: color.RGBA{0, 0, 0, 255},       // f => black,
}

// Gets a MD5 hash from the given string limited to exactly 64 chars.
//
func getMD5(s string) string {
	length := len(s)
	if length < 16 {
		s = s + enlarger[0:16-length]
	}
	hasher := md5.New()
	return fmt.Sprintf("%x", hasher.Sum([]byte(s)))
}

// Generates a image.RGBA of colors given a hexadecimal string. It draws solid rectangles
// instead of iterate over all the pixels of the image.
// Due to this, it always performs the same number of painting operations (64) and
// times are much more constant.
//
// Time required for generate an image 64x64:  		 0.005s
// Time required for generate an image 4096x4096:  0.613s
//
func buildImage(hash string, scale int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, avatarSide*scale, avatarSide*scale))

	for x := 0; x < avatarSide; x++ {
		for y := 0; y < avatarSide; y++ {
			color := palette[hash[x*avatarSide+y]]
			startPoint := image.Point{x * scale, y * scale}
			endPoint := image.Point{x*scale + scale, y*scale + scale}
			rectangle := image.Rectangle{startPoint, endPoint}
			draw.Draw(img, rectangle, &image.Uniform{color}, image.ZP, draw.Src)
		}
	}

	return img
}

// Exports the given image in the given path.
//
func exportImage(img *image.RGBA, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer file.Close()

	extension := extensionRegexp.FindString(filepath)

	if extension == ".png" {
		png.Encode(file, img)
	} else if extension == ".jpg" || extension == ".jpeg" {
		jpeg.Encode(file, img, &jpeg.Options{95})
	} else {
		fmt.Println("Invalid extension")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Avatarme"
	app.Usage = "Generates an unique avatar for the given string"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "output, o", Value: "output.png", Usage: "path of the output file"},
		cli.IntFlag{Name: "size, s", Value: 256, Usage: "side length of the generated image (in px). Will be ronded to a multiple of 8"},
	}
	app.Action = func(c *cli.Context) {
		text := c.Args()[0]
		hash := getMD5(text)
		img := buildImage(hash, c.Int("size")/8)
		exportImage(img, c.String("output"))
	}

	app.Run(os.Args)
}
