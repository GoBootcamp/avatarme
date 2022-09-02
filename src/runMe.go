package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	// make a hash of the personal info
	userHash := getHash()
	fmt.Println("The hashed value is: ", userHash)

	//hashSize := len(userHash) // 88 for sha 512
	bw := []color.Color{color.Black, color.White}
	outputImage := image.NewPaletted(
		image.Rect(0, 0, 100, 200),
		bw,
	)

	writeImage(outputImage)
	//outputPath := time.Now().Format("01-02-2006 15:04:05") + ".png"

	// TODO:
	// make an image of the hash of the IP address
	// save to file?
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
func getHash() string {
	// get a string (eg IP address, email)
	fmt.Println("Enter Your Personal Information: ")
	var personalInfo string
	fmt.Scanln(&personalInfo)

	hasher := sha512.New()
	bv := []byte(personalInfo)
	hasher.Write(bv)

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

/**
import (
    "image"
)
func main() {
    myImg := image.NewRGBA(image.Rect(0, 0, 12, 6))
        out, err := os.Create("cat.png")
        png.Encode(out, myImg)
        out.Close()
}
*/
