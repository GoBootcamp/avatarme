package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/ymakhloufi/ydenticon/ydenticon"
)

// ToDo: make output into file optional via flag and print into stdout by default to allow piping
// ToDo: for optional filename, prompt if file exists, also provide force-flag for overwriting existing file
// ToDo: allow to control complexity and image width vie CLI arg
// ToDo: read up on docs and how to provide CLI help
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

	fmt.Printf("Image was written to %s\n", outputFilePtr.Name())
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
  ydenticon <identifier> [--complexity=<level>] [width=<widthInPx>] [(--output=<filePathAndName> | --outputOverwrite=<filePathAndName>)]
  ydenticon -h | --help
  ydenticon --version

Options:
  -h --help                                                 Show this screen.
  -v --version                                              Show version.
  -o=<filePathAndName> --output=<filePathAndName>           Path and file name of output file [default: <STDOUT>].
  -O=<filePathAndName> --outputOverwrite=<filePathAndName>  Path and file name of output file (overwrite if exists).
  -c=<level> --complexity=<level>                           Result's level of complexity ( 1 | 2 | 3 | 4 | 5 ) [default: 3].
  -w=<widthInPx> --width=<widthInPx>                        Result image's width in pixels [default: 200].`

	arguments, _ := docopt.ParseDoc(usage)

	identifier, _ = arguments.String("<identifier>")

	width, _ := arguments.Int("--width")
	widthInPx = uint(width)

	complexityInt, _ := arguments.Int("--complexity")
	complexityLevel, err := ydenticon.GetComplexityLevel(complexityInt)
	ExitIfErr(err)

	outputPath, _ := arguments.String("--output")
	if outputPath == "<STDOUT>" {
		output = os.Stdout
	} else {
		// ToDo: Find out why I can't assign 'os.Create' directly to 'output' (something about scoping of ':=' vs '=' ?)
		overwriteExistingFiles, _ := arguments.Bool("overwrite")
		if stat, err := os.Stat(outputPath); os.IsNotExist(err) || (overwriteExistingFiles && !stat.IsDir()) {
			file, err := os.Create(outputPath)
			ExitIfErr(err)
			output = file
		} else if stat.IsDir() {
			ExitIfErr(fmt.Errorf("%s exists but. To overwrite use --overwrite flag", outputPath))
		} else if err != nil {

		} else {
			ExitIfErr(fmt.Errorf("file %s exists. To overwrite use --overwrite flag", outputPath))
		}
	}

	return
}

func ExitIfErr(err error) {
	// ToDo: add test
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
