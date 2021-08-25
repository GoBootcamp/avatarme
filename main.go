package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kennygrant/sanitize"

	"github.com/ymakhloufi/avatarme/identicon"
)

// ToDo: make output into file optional via flag and print into stdout by default to allow piping
// ToDo: for optional filename, prompt if file exists, also provide force-flag for overwriting existing file
// ToDo: allow to control complexity and image width vie CLI arg
// ToDo: read up on docs and how to provide CLI help
// ToDo: read up on info/debug logging
// ToDo: read up on debugging/stepping-though code
func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to pass the textual identifier as an argument.")
		os.Exit(1)
	}
	identifier := os.Args[1]
	identiconObj := identicon.New(identifier)
	fileName := sanitize.BaseName(identifier)
	outputPath, err := filepath.Abs(filepath.FromSlash("output/" + fileName + ".png"))
	exitIfErr(err)

	err = identiconObj.SavePngToDisk(outputPath, identicon.ComplexityLevelMedium, 200)
	exitIfErr(err)

	fmt.Printf("Identicon was written to %s\n", outputPath)
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
