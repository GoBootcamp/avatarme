package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/ymakhloufi/ydenticon/ydenticon"
)

// ToDo: read up on info/debug logging
// ToDo: read up on debugging/stepping-though code
func main() {
	identifier, outputFilePtr, complexityLevel, widthInPx := getCliArgs()
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			if _, printErr := fmt.Fprintln(os.Stderr, err); printErr != nil {
				fmt.Println("error printing to STDERR: ", printErr, err)
			}
			os.Exit(1)
		}
	}(outputFilePtr)

	ydenticonObj := ydenticon.New(identifier)
	err := ydenticonObj.SavePngToDisk(outputFilePtr, complexityLevel, widthInPx)
	ExitIfErr(err)
}

// ToDo: add test
func getCliArgs() (
	identifier string,
	output *os.File,
	complexityLevel ydenticon.ComplexityLevel,
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

	arguments, err := docopt.ParseDoc(usage)
	ExitIfErr(err)
	// panic(fmt.Sprintf("%v", arguments))

	identifier, err = arguments.String("<identifier>")
	ExitIfErr(err)

	width, err := arguments.Int("--width")
	ExitIfErr(err)
	widthInPx = uint(width)

	complexityInt, err := arguments.Int("--complexity")
	ExitIfErr(err)
	complexityLevel, err = ydenticon.GetComplexityLevel(complexityInt)
	ExitIfErr(err)

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

	return
}

func ExitIfErr(err error) {
	// ToDo: add test
	if err != nil {
		println(fmt.Sprintf("%v\n", err))
		os.Exit(1)
	}
}
