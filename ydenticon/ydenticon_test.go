package ydenticon

import (
	"image/color"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ToDo: read up on mocking, write unit test(s) for drawSquares()

func TestByteSliceToBoolSlice(t *testing.T) {
	testData := []struct {
		input  []byte
		expect []bool
	}{
		{
			input: []byte{255, 0},
			expect: []bool{
				true, true, true, true, true, true, true, true,
				false, false, false, false, false, false, false, false,
			},
		},
		{
			input: []byte{85, 170},
			expect: []bool{
				true, false, true, false, true, false, true, false,
				false, true, false, true, false, true, false, true,
			},
		},
	}

	for _, data := range testData {
		result := byteSliceToBoolSlice(data.input)
		assert.Equal(t, len(data.expect), len(result))
		assert.Equal(t, data.expect, result)
	}
}

func TestGetColorsFromBytesReturnsSameValuesAsForegroundColor(t *testing.T) {
	// 20 random sets of test data (0-255 for each respective RGB value)
	for i := 0; i < 20; i++ {
		r, g, b := byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))
		expectedColor := color.RGBA{r, g, b, 0xff}
		result, _ := getColorsFromBytes(r, g, b)
		assert.Equal(t, expectedColor, result)
	}
}

func TestGetColorsFromBytesReturnsCorrectBackgroundColor(t *testing.T) {
	testData := []struct {
		input  struct{ r, g, b byte }
		expect color.RGBA
	}{
		{
			// On the verge, but still black
			input:  struct{ r, g, b byte }{168, 221, 55},
			expect: color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
		},
		{
			// On the verge, but still white
			input:  struct{ r, g, b byte }{168, 220, 55},
			expect: color.RGBA{R: 255, G: 255, B: 255, A: 0xff},
		},
		{
			input:  struct{ r, g, b byte }{255, 255, 255},
			expect: color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
		},
		{
			input:  struct{ r, g, b byte }{0, 0, 0},
			expect: color.RGBA{R: 255, G: 255, B: 255, A: 0xff},
		},
	}
	for _, data := range testData {
		_, result := getColorsFromBytes(data.input.r, data.input.g, data.input.b)
		assert.Equal(t, data.expect, result)
	}
}
