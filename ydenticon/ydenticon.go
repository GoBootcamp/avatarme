package ydenticon

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
	"time"

	"github.com/nfnt/resize"
)

type ArtifactDimension uint8

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
	ComplexityLevelLowest ArtifactDimension = 5 + (1 << iota)
	ComplexityLevelLow
	ComplexityLevelMedium
	ComplexityLevelHigh
	ComplexityLevelUltra
)

func GetComplexityLevel(level int) (ArtifactDimension, error) {
	levels := [...]ArtifactDimension{
		ComplexityLevelLowest,
		ComplexityLevelLow,
		ComplexityLevelMedium,
		ComplexityLevelHigh,
		ComplexityLevelUltra,
	}
	index := level - 1
	if index > len(levels)-1 || index < 0 {
		return 0, fmt.Errorf("the chosen complexity level %d does not exist, Try 1-5", level)
	}
	return levels[index], nil
}

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
		backgroundColor = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	} else {
		backgroundColor = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	}

	foregroundColor = color.RGBA{R: r, G: g, B: b, A: 0xff}
	return foregroundColor, backgroundColor
}

func (i Ydenticon) Make(dimension ArtifactDimension, widthInPx uint) (*image.RGBA, error) {
	cols, rows := int(dimension), int(dimension) // slightly redundant, but we might want non-square-shaped images later.
	artifactSizeInPx := int(math.Ceil(float64(widthInPx) / float64(dimension)))
	canvasPtr := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: artifactSizeInPx * cols, Y: artifactSizeInPx * rows},
		},
	)
	draw.Draw(canvasPtr, canvasPtr.Bounds(), &image.Uniform{C: i.backgroundColor}, image.Point{X: 0, Y: 0}, draw.Src)
	drawSquares(canvasPtr, i.canvasBitSlice, i.foregroundColor, cols, rows, artifactSizeInPx)

	if resizedCanvas, ok := resize.Resize(widthInPx, widthInPx, canvasPtr, resize.Lanczos3).(*image.RGBA); ok == true {
		canvasPtr = resizedCanvas
	} else {
		return canvasPtr, fmt.Errorf("failed extracting RGBA type from resized image: %v", resizedCanvas)
	}

	return canvasPtr, nil
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
func drawSquares(canvas *image.RGBA, bits []bool, color color.RGBA, cols, rows int, artifactSizeInPx int) {
	bitIndex := 0
	// only go through half of the columns and mirror the squares onto the second half
	middleColumn := int(math.Ceil(float64(cols) / 2))

	var wg sync.WaitGroup
	for x := 0; x < middleColumn; x++ {
		for y := 0; y < rows; y++ {
			if !bits[bitIndex] {
				bitIndex++
				continue
			}
			wg.Add(1)
			go drawOneSquare(&wg, canvas, x, y, artifactSizeInPx, middleColumn, cols, color)

			bitIndex++
		}
	}

	wg.Wait()
}

func drawOneSquare(
	wg *sync.WaitGroup,
	canvas *image.RGBA,
	x, y, artifactSizeInPx, middleColumn, cols int,
	color color.RGBA,
) {
	defer wg.Done()
	location := image.Point{X: x * artifactSizeInPx, Y: y * artifactSizeInPx}
	rect := image.Rectangle{
		Min: location,
		Max: image.Point{location.X + artifactSizeInPx, location.Y + artifactSizeInPx},
	}

	draw.Draw(canvas, rect, &image.Uniform{C: color}, image.Point{X: 0, Y: 0}, draw.Src)

	// Don't mirror anything if we are working on the middle column if there is an uneven number of columns,
	// since it would just copy the artifact onto itself (same coords) which is redundant work.
	if cols%2 == 0 || x < middleColumn {
		mirroredRect := image.Rectangle{
			Min: image.Point{X: artifactSizeInPx*(cols-1) - location.X, Y: location.Y},
			Max: image.Point{X: artifactSizeInPx*cols - location.X, Y: location.Y + artifactSizeInPx},
		}
		draw.Draw(canvas, mirroredRect, &image.Uniform{C: color}, image.Point{X: 0, Y: 0}, draw.Src)
	}
	// Time difference gets visible when you run this sequentially instead of inside a goroutine.
	// This sleep should obviously be removed in a production system!
	time.Sleep(10 * time.Millisecond)
}
