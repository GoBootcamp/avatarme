package ydenticon

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

type ComplexityLevel uint8

const (
	// ComplexityLevelLowest ComplexityLevelLow ComplexityLevelMedium ComplexityLevelHigh ComplexityLevelUltra
	/* Used as dimensions for the image. SHA256 gives us 256 bits for uniqueness.
	 * After determining the color (3 bytes = 24 bits) we have 232 bits left.
	 * ComplexityLevelUltra uses 231 bits, so we would have 1 bit left to spare.
	 * Lowest: 6x6   = Ceil(6/2) * 6   = 18 bits
	 * Low:    7x7   = Ceil(7/2) * 7   = 28 bits
	 * Medium: 9x9   = Ceil(9/2) * 9   = 81 bits
	 * High:   13x13 = Ceil(13/2) * 13 = 91 bits
	 * Ultra:  21x21 = Ceil(21/2) * 21 = 231 bits
	 */
	ComplexityLevelLowest ComplexityLevel = 5 + (1 << iota)
	ComplexityLevelLow
	ComplexityLevelMedium
	ComplexityLevelHigh
	ComplexityLevelUltra
)

type Ydenticon struct {
	Identifier                       string
	hashByteArray                    [32]byte
	canvasBitSlice                   []bool
	foregroundColor, backgroundColor color.RGBA
}

func New(identifier string) *Ydenticon {
	hashByteArray := sha256.Sum256([]byte(identifier))
	foregroundColor, backgroundColor := getColorsFromBytes(hashByteArray[0], hashByteArray[1], hashByteArray[2])

	return &Ydenticon{
		Identifier:      identifier,
		hashByteArray:   hashByteArray,
		canvasBitSlice:  byteSliceToBoolSlice(hashByteArray[3:]),
		foregroundColor: foregroundColor,
		backgroundColor: backgroundColor,
	}
}

func getColorsFromBytes(r, g, b byte) (foregroundColor color.RGBA, backgroundColor color.RGBA) {
	// http://stackoverflow.com/a/3943023/112731 (contrast-background black or white)
	if (0.299*float32(int(r)) + 0.587*float32(int(g)) + 0.114*float32(int(b))) > 186 {
		backgroundColor = color.RGBA{0, 0, 0, 0xff}
	} else {
		backgroundColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
	}

	foregroundColor = color.RGBA{r, g, b, 0xff}
	return foregroundColor, backgroundColor
}

func (i Ydenticon) SavePngToDisk(filePath string, level ComplexityLevel, widthInPx uint) error {
	cols, rows := int(level), int(level) // slightly redundant, but we might want non-square-shaped images later.
	elemSizeInPx := int(math.Ceil(float64(level) / float64(widthInPx)))
	canvas := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{elemSizeInPx * cols, elemSizeInPx * rows},
		},
	)
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{i.backgroundColor}, image.Point{0, 0}, draw.Src)
	drawSquares(canvas, i.canvasBitSlice, i.foregroundColor, cols, rows, elemSizeInPx)

	// Todo: ensure exact width by re-scaling png before writing to disk

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error opening/creating png file on disk %s: %v", filePath, err)
	}

	if err = png.Encode(file, canvas); err != nil {
		return fmt.Errorf("error encoding png file %s: %v", filePath, err)
	}

	return nil
}

func byteSliceToBoolSlice(bytes []byte) []bool {
	var result []bool
	for _, myByte := range bytes {
		for i := 0; i <= 7; i++ {
			result = append(result, myByte&(1<<i) > 0)
		}
	}

	return result
}

// This only fills the left half of the canvas and mirrors it to the right half
// which reduces the bits needed (and thus uniqueness-per-surface-area)
// but makes it significantly more human-recognizable.
func drawSquares(canvas *image.RGBA, bits []bool, color color.RGBA, cols, rows int, elemSizeInPx int) {
	bitIndex := 0
	// only go through half of the columns and mirror the squares onto the second half
	middleColumn := int(math.Ceil(float64(cols) / 2))

	// ToDo: parallelize? Need a mutex on the canvas object? Can solve with GoRoutines even though resources are shared?
	for x := 0; x < middleColumn; x++ {
		for y := 0; y < rows; y++ {
			if !bits[bitIndex] {
				bitIndex++
				continue
			}
			location := image.Point{x * elemSizeInPx, y * elemSizeInPx}
			rect := image.Rectangle{
				location,
				image.Point{location.X + elemSizeInPx, location.Y + elemSizeInPx},
			}
			draw.Draw(canvas, rect, &image.Uniform{color}, image.Point{0, 0}, draw.Src)

			// Don't mirror anything if we are working on the middle column if there is an uneven number of columns,
			// since it would just copy the element onto itself (same coords) which is redundant work.
			if cols%2 == 0 || x < middleColumn {
				mirroredRect := image.Rectangle{
					image.Point{elemSizeInPx*(cols-1) - location.X, location.Y},
					image.Point{elemSizeInPx*cols - location.X, location.Y + elemSizeInPx},
				}
				draw.Draw(canvas, mirroredRect, &image.Uniform{color}, image.Point{0, 0}, draw.Src)
			}
			bitIndex++
		}
	}
}
