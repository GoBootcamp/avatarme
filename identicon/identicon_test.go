package identicon

import (
	"image/color"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// ToDo: read up on mocking, write unit test(s) for drawSquares()
// ToDo: read up on assertion frameworks, rewrite tests with such to make more concise/legible.

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
		if len(result) != len(data.expect) {
			t.Errorf("length of input was %d, expected %d", len(result), len(data.expect))
		}

		for k, v := range data.expect {
			if v != result[k] {
				t.Errorf("%v, was expected to be %v", result, data.expect)
			}
		}
	}
}

func TestGetColorsFromBytesReturnsSameValuesAsForegroundColor(t *testing.T) {
	// 20 random sets of test data (0-255 for each respective RGB value)
	for i := 0; i < 20; i++ {
		r, g, b := byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))
		expectedColor := color.RGBA{r, g, b, 0xff}
		result, _ := getColorsFromBytes(r, g, b)
		if !cmp.Equal(result, expectedColor) {
			t.Errorf("Foreground color was %#v, expected %#v", result, expectedColor)
		}
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
		if !cmp.Equal(data.expect, result) {
			t.Errorf("Foreground color was %#v, expected %#v", result, data.expect)
		}
	}
}
