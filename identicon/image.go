package identicon

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

const rows, cols = 8, 8
const elemSizeInPx = 20

type Identicon struct {
	Identifier, HashString string
	HashByteArray          [32]byte
}

func New(identifier string) *Identicon {
	hashByteArray := sha256.Sum256([]byte(identifier))
	hashByteSlice := hashByteArray[:]
	hashString := hex.EncodeToString(hashByteSlice)
	return &Identicon{
		Identifier:    identifier,
		HashByteArray: hashByteArray,
		HashString:    hashString,
	}
}

func (i Identicon) MakeImageFromHash(filePath string) {
	var colorBytes [3]byte
	copy(colorBytes[:], i.HashByteArray[:3])
	foregroundColor, backgroundColor := getColors(i.HashByteArray[0], i.HashByteArray[1], i.HashByteArray[2])

	canvas := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{elemSizeInPx * cols, elemSizeInPx * rows}})
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{backgroundColor}, image.Point{0, 0}, draw.Src)

	// need one bit per drawn element, skip the first 3 bytes because they were already used for color above
	neededBytes := int(math.Ceil(cols * rows / 8))
	sourceBytes := i.HashByteArray[3 : 3+neededBytes]
	drawSquares(canvas, byteSliceToBoolSlice(sourceBytes), foregroundColor)

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	if err = png.Encode(file, canvas); err != nil {
		panic(err)
	}
}

func getColors(r, g, b byte) (color.RGBA, color.RGBA) {
	foregroundColor := color.RGBA{r, g, b, 0xff}
	var backgroundColor color.RGBA

	// http://stackoverflow.com/a/3943023/112731 (contrast-background black or white)
	if (0.299*float32(int(r)) + 0.587*float32(int(g)) + 0.114*float32(int(b))) > 186 {
		backgroundColor = color.RGBA{0, 0, 0, 0xff}
	} else {
		backgroundColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
	}

	return foregroundColor, backgroundColor
}

func byteSliceToBoolSlice(bytes []byte) []bool {
	var result []bool
	for _, myByte := range bytes {
		for i := 7; i >= 0; i-- {
			result = append(result, myByte&(1<<i) > 0)
		}
	}

	return result
}

func drawSquares(canvas *image.RGBA, bits []bool, color color.RGBA) {
	bitIndex := 0
	for x := 0; x < cols; x++ {
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
			mirroredRect := image.Rectangle{
				image.Point{elemSizeInPx*(cols-1) - location.X, location.Y},
				image.Point{elemSizeInPx*cols - location.X, location.Y + elemSizeInPx},
			}
			draw.Draw(canvas, rect, &image.Uniform{color}, image.Point{0, 0}, draw.Src)
			draw.Draw(canvas, mirroredRect, &image.Uniform{color}, image.Point{0, 0}, draw.Src)
			bitIndex++
		}
	}
}
