package pixelated

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

const (
	avatarSide      int    = 8
	signatureLength int    = avatarSide * avatarSide
	enlarger        string = "md5-enlarger"
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

// Generates a MD5 hash from the given string limited to exactly 64 chars.
// If the length of the string is not enough to generate a valid hash, appends some content
// at the end.
//
func getMD5(s string) string {
	length := len(s)
	if length < 16 {
		s = s + enlarger[0:16-length]
	}
	hasher := md5.New()
	return fmt.Sprintf("%x", hasher.Sum([]byte(s)))
}

// Generates a image.RGBA of colors given a hexadecimal string.
//
func BuildImage(text string, size int) *image.RGBA {
	hash := getMD5(text)
	scale := size / avatarSide

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
