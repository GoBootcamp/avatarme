package encoder

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"regexp"
)

var extensionRegexp *regexp.Regexp = regexp.MustCompile(`\.(png|jpe?g)$`)

// Exports the given image in the given path.
//
func ExportImage(img *image.RGBA, filepath string) {
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
