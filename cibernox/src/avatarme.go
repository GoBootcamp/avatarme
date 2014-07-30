package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const (
	enlarger        string = "md5-enlarger"
	avatarSide      int    = 8
	signatureLength int    = avatarSide * avatarSide
	scale           int    = 512
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
	hasher := md5.New()
	byteSignature := hasher.Sum([]byte(s + enlarger))
	chunkedSignature := byteSignature[0 : signatureLength/2]
	return fmt.Sprintf("%x", chunkedSignature)
}

// Generates a image.RGBA of colors given a hexadecimal string.
//
// Time required for compile and generate an image 64x64:  		 0.175s
// Time required for compile and generate an image 4096x4096:  1.644s
//
func buildImagePixelByPixel(hash string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, avatarSide*scale, avatarSide*scale))

	for x := 0; x < avatarSide; x++ {
		for y := 0; y < avatarSide; y++ {
			fillPixel(img, x*scale, y*scale, palette[hash[x*avatarSide+y]])
		}
	}

	return img
}

// Generates a image.RGBA of colors given a hexadecimal string. It draws solid rectangles
// instead of iterate over all the pixels of the image.
// Due to this, it always performs the same number of painting operations (64) and
// times are much more constant.
//
// Time required for compile and generate an image 64x64:  		 0.176s
// Time required for compile and generate an image 4096x4096:  0.784s
//
func buildImageDrawing(hash string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, avatarSide*scale, avatarSide*scale))

	for x := 0; x < avatarSide; x++ {
		for y := 0; y < avatarSide; y++ {
			color := palette[hash[x*avatarSide+y]]
			startPoint := image.Point{x*scale, y*scale}
			endPoint := image.Point{x*scale + scale, y*scale + scale}
			rectangle := image.Rectangle{startPoint, endPoint}
			draw.Draw(img, rectangle, &image.Uniform{color}, image.ZP, draw.Src)
		}
	}

	return img
}

// Fill the virtual pixel with the given coordinates with the given color.
//
// Images always are a 8x8 virtual image, but the output might be bigger.
// Per example, if we want a 512x512 avatar, we set the `scale` constant to 64.
// That way each virtual pixel will be a square of 32x32, resulting in a 512x512 image.
//
func fillPixel(img *image.RGBA, x, y int, color color.RGBA) {
	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			img.Set(x+i, y+j, color)
		}
	}
}

// Exports the given image as a PNG file with the name `output.png`
//
func exportImage(img *image.RGBA) {
	file, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer file.Close()

	png.Encode(file, img)
}

func main() {
	text := os.Args[1]
	hash := getMD5(text)
	// img := buildImagePixelByPixel(hash)
	img := buildImageDrawing(hash)
	exportImage(img)
}
