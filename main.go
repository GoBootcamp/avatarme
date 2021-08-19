package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ymakhloufi/avatarme/identicon"
)

// ToDo: read up on docs and how to provide CLI help
// ToDo: allow to control
func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to pass the textual identifier as an argument.")
		os.Exit(1)
	}
	identiconObj := identicon.New(os.Args[1])

	// ToDo: sanitize os.Args[1]
	relativeFilePath := filepath.FromSlash("output/" + os.Args[1] + ".png")
	if err := identiconObj.SavePngToDisk(relativeFilePath, identicon.ComplexityLevelMedium, 200); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cwd, _ := os.Getwd() // if error happens during Getwd(), we ignore it and just print the relative path
	fullPath := filepath.Join(cwd, relativeFilePath)
	fmt.Printf("Identicon was written to %s\n", fullPath)
}
