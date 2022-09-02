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

// main
//
// from the src directory, type 'go run .'
// when prompted, enter the unique personal information (email, IP address)
// files are outputted in the same working directory, under the name identicon.png
func main() {
	// make a hash of the personal info
	userHash := getHash()

	// create the image
	outputImage := createImage(userHash)

	// for the purposes of this PR, write image to the disk
	writeImage(outputImage)
}

// createImage
//
// Given a user hash, creates an associated image
func createImage(userHash []byte) *image.Paletted {
	// 64 for sha 512.
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
	return outputImage
}

// writeImage
//
// Writes out the image to a hard-coded location on disk, named identicon.png
func writeImage(outputImage *image.Paletted) {

	outputPath := "identicon.png"
	out, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Could not create image")
		return
	}
	png.Encode(out, outputImage)
	defer out.Close()
	fmt.Println("identicon written to ", outputPath)
}

// getHash
//
// Returns the hash from a user
func getHash() []byte {
	// get a string (eg IP address, email)
	fmt.Println("Enter Your Personal Information: ")
	var personalInfo string
	fmt.Scanln(&personalInfo)

	hasher := sha512.New()
	bv := []byte(personalInfo)
	hasher.Write(bv)
	userHash := hasher.Sum(nil)
	fmt.Println("The hashed value is: ", base64.URLEncoding.EncodeToString(userHash))
	return userHash
}
