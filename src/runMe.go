package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	// make a hash of the personal info
	userHash := getHash()
	fmt.Println("The hashed value is: ", base64.URLEncoding.EncodeToString(userHash))

	// 64 for sha 256
	hashSize := len(userHash)
	gridWidth := 8
	scale := 40
	bw := []color.Color{color.Black, color.White}
	outputImage := image.NewPaletted(
		image.Rect(0, 0, gridWidth*scale, gridWidth*scale),
		bw,
	)

	for i := 0; i < hashSize; i++ {
		hashValue := userHash[i]
		myColor := bw[hashValue%2]

		// row = index % width
		// i ==10 corresponds to row 1, cell 2
		x := (i / gridWidth) * scale
		y := (i % gridWidth) * scale

		start := image.Point{x, y}
		end := image.Point{x + scale, y + scale}
		rectangle := image.Rectangle{start, end}
		draw.Draw(outputImage, rectangle, &image.Uniform{myColor}, image.Point{}, draw.Src)
	}

	writeImage(outputImage)

	// TODO:
	// make an image of the hash of the IP address
}

func writeImage(outputImage *image.Paletted) {

	outputPath := "identicon.png"
	out, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Could not create image")
		return
	}
	png.Encode(out, outputImage)
	out.Close()
	fmt.Println("identicon written to ", outputPath)
}
func getHash() []byte {
	// get a string (eg IP address, email)
	fmt.Println("Enter Your Personal Information: ")
	var personalInfo string
	fmt.Scanln(&personalInfo)

	hasher := sha512.New()
	bv := []byte(personalInfo)
	hasher.Write(bv)

	//sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	//return sha
	return hasher.Sum(nil)
}
