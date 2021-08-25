package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kennygrant/sanitize"

	"github.com/ymakhloufi/ydenticon/ydenticon"
)

// ToDo: make output into file optional via flag and print into stdout by default to allow piping
// ToDo: for optional filename, prompt if file exists, also provide force-flag for overwriting existing file
// ToDo: allow to control complexity and image width vie CLI arg
// ToDo: read up on docs and how to provide CLI help
// ToDo: read up on info/debug logging
// ToDo: read up on debugging/stepping-though code
func main() {
	identifier, absOutputPath, complexityLevel, widthInPx := getCliArgs()

	ydenticonObj := ydenticon.New(identifier)
	err := ydenticonObj.SavePngToDisk(absOutputPath, complexityLevel, widthInPx)
	exitIfErr(err)

	fmt.Printf("Image was written to %s\n", absOutputPath)
}

// ToDo: add test
func getCliArgs() (
	identifier string,
	absOutputPath string,
	complexityLevel ydenticon.ComplexityLevel,
	widthInPx uint,
) {
	// ToDo: do away with inbuilt flags. Find package to cleanly manage flexible arg-order and required args/flags
	// ToDo: docopt vs. jessevdk/go-flags

	flag.StringVar(&identifier, "id", "", "path and filename to output file")
	fileName := *flag.String("output", sanitize.BaseName(identifier), "path and filename to output file")
	absOutputPath, err := filepath.Abs(filepath.FromSlash("output/" + fileName + ".png"))
	// // ToDo: add flag to create directory if not exists
	fmt.Println(absOutputPath)
	exitIfErr(err)

	complexityLevel = ydenticon.ComplexityLevelMedium

	widthInPx = 200

	flag.Parse()

	tail := flag.Args()
	if len(tail) > 1 {
		exitIfErr(fmt.Errorf("too many args after flag. Only one is accepted, but received %v", tail))
		// ToDo: print proper CLI usage or help-page instead of this error
	}

	return
}

// ToDo: add test
func exitIfErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
