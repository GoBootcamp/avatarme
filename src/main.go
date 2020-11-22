package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/png"
	"log"
	"os"
)

const hashLength = 40

func getHash(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

func switchIndexDirection() func() int {
	idx := 0
	currentSign := "+"
	positiveSign := "+"
	negativeSign := "-"

	return func() int {
		if idx == 0 {
			currentSign = positiveSign
		} else if idx == hashLength-1 {
			currentSign = negativeSign
		}
		if currentSign == positiveSign {
			idx++
			return idx
		}
		idx--
		return idx
	}
}

func generateImage(s string, size, sideBlocks int) *image.RGBA {
	scale := size / sideBlocks
	img := image.NewRGBA(image.Rect(0, 0, sideBlocks*scale, sideBlocks*scale))

	idx := switchIndexDirection()
	for x := 0; x < sideBlocks; x++ {
		for y := 0; y < sideBlocks; y++ {
			col := palette.Plan9[s[idx()]]
			startPoint := image.Point{x * scale, y * scale}
			endPoint := image.Point{x*scale + scale, y*scale + scale}
			rectangle := image.Rectangle{startPoint, endPoint}
			draw.Draw(img, rectangle, &image.Uniform{col}, image.Point{}, draw.Src)
		}
	}

	return img
}

func saveImage(img *image.RGBA, filename string) {

	f, err := os.Create(fmt.Sprintf("%s.png", filename))
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

// Run it as: go run ./main.go -input=Awesome -size 256 -output image -blocks 8
func main() {
	inputString := flag.String("input", "Cool", "value to be hashed and generate identicon with")
	outputFile := flag.String("output", "image", "image file name")
	imageSize := flag.Int("size", 256, "image size")
	blocks := flag.Int("blocks", 8, "number of side blocks")

	flag.Parse()

	hash := getHash(*inputString)
	fmt.Println(hash)

	img := generateImage(hash, *imageSize, *blocks)
	saveImage(img, *outputFile)
}
