package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"

	"github.com/ymakhloufi/ydenticon/ydenticon"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: time.RFC3339, FullTimestamp: true})

	// For SQL-style date-time, use this EXACTLY (cannot change digits): TimestampFormat: "2006-01-02 15:04:05"
	// Explanation: https://golang.org/src/time/format.go
}

// ToDo: read up on debugging/stepping-though code
func main() {
	identifier, outputFilePtr, artifactDimension, widthInPx := getCliArgs()
	defer func(file *os.File) {
		log.Debugf("Closing file %s", file.Name())
		err := file.Close()
		if err != nil {
			ExitIfErr(fmt.Errorf("Failed to close file %s", file.Name()))
		}
	}(outputFilePtr)

	log.Infof("Identifier:        %s", identifier)
	log.Infof("Output File:       %s", outputFilePtr.Name())
	log.Infof("Complexity:        %d x %d artifacts", artifactDimension, artifactDimension)
	log.Infof("Image Dimensions:  %d x %d Pixels", widthInPx, widthInPx)

	log.Debug("Creating Ydenticon object")
	ydenticonObj := ydenticon.New(identifier)

	log.Debugf("Writing output to %s", outputFilePtr.Name())
	img, err := ydenticonObj.Make(artifactDimension, widthInPx)
	if err != nil {
		ExitIfErr(fmt.Errorf("fialed to make identicon: %v", err))
	}

	if err = png.Encode(outputFilePtr, img); err != nil {
		ExitIfErr(fmt.Errorf("error encoding png file %s: %v", outputFilePtr.Name(), err))
	}
}

// ToDo: add test
func getCliArgs() (
	identifier string,
	output *os.File,
	artifactDimension ydenticon.ArtifactDimension,
	widthInPx uint,
) {
	usage := `Ydention generates a unique avatar ("Identicon") based on a given unique string identifier.

Usage:
  ydenticon <identifier> [--complexity=<level>] [--width=<widthInPx>] [(--output=<filePathAndName> [--overwriteExistingFile])]
  ydenticon -h | --help
  ydenticon --version

Options:
  -h --help                                Show this screen.
  -v --version                             Show version.
  -o, --output=<filePathAndName>           Path and file name of output file [default: <STDOUT>].
  -f, --overwriteExistingFile              Overwrite target file if exists [default: false]
  -c, --complexity=<level>                 Result's level of complexity ( 1 | 2 | 3 | 4 | 5 ) [default: 3].
  -w, --width=<widthInPx>                  Result image's width in pixels [default: 200].`

	// ToDo: figure out a better way to deal with many errors like this.

	log.Debug("Parsing DocOpt string")
	arguments, err := docopt.ParseDoc(usage)
	ExitIfErr(err)
	log.Debugf("Found These CLI args:\n%v", arguments)

	log.Debugf("Parsing 'identifier' from CLI args")
	identifier, err = arguments.String("<identifier>")
	ExitIfErr(err)
	log.Debugf("Found 'identifier' %s", identifier)

	log.Debugf("Parsing 'width' from CLI args")
	width, err := arguments.Int("--width")
	ExitIfErr(err)
	widthInPx = uint(width)
	log.Debugf("Found 'width' %d", widthInPx)

	log.Debugf("Parsing 'complexity' from CLI args")
	complexityInt, err := arguments.Int("--complexity")
	ExitIfErr(err)
	log.Debugf("Found 'complexity' %s", complexityInt)

	log.Debug("Converting complexity level to image dimensions")
	artifactDimension, err = ydenticon.GetComplexityLevel(complexityInt)
	ExitIfErr(err)
	log.Debugf(
		"Computed Artifact Dimensions from Complexity Level %d: %dx%d artifacts",
		complexityInt,
		artifactDimension,
		artifactDimension,
	)

	if widthInPx < uint(artifactDimension) {
		ExitIfErr(
			fmt.Errorf(
				"image width %d is too small. It must be at least %d pixels for complexity level %d",
				widthInPx,
				artifactDimension,
				complexityInt,
			),
		)
	}

	log.Debugf("Parsing 'output' from CLI args")
	outputPath, err := arguments.String("--output")
	ExitIfErr(err)
	if outputPath == "<STDOUT>" {
		output = os.Stdout
	} else {
		// ToDo: Find out why I can't assign 'os.Create' directly to 'output' (something about scoping of ':=' vs '=' ?)
		overwriteExistingFiles, err := arguments.Bool("--overwriteExistingFile")
		ExitIfErr(err)
		if os.IsNotExist(err) || overwriteExistingFiles {
			output, err = os.Create(outputPath) // ToDo: read up on `go vet` to forbid var-shadowing
			ExitIfErr(err)
		} else if err != nil {
			ExitIfErr(err) // other error resulting from os.stat()
		} else {
			ExitIfErr(fmt.Errorf("file %s exists. To overwrite use the --overwrite flag", outputPath))
		}
	}
	log.Debugf("Found 'output' %s", output.Name())

	return
}

func ExitIfErr(err error) {
	// ToDo: add test
	// https://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go
	if err != nil {
		log.Fatalf("%v\n", err)
		os.Exit(1)
	}
}
