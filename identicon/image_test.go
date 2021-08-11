package identicon

import "testing"

func TestByteSliceToBoolSlice(t *testing.T) {
	byteSlice := []byte{255, 0, 85, 170}
	expectedBoolSlice := []bool{
		true, true, true, true, true, true, true, true,
		false, false, false, false, false, false, false, false,
		true, false, true, false, true, false, true, false,
		false, true, false, true, false, true, false, true,
	}
	result := byteSliceToBoolSlice(byteSlice)
	if len(result) != 32 {
		t.Errorf("length of input was %d, expected %d", len(result), 32)
	}

	for k, v := range expectedBoolSlice {
		if v != result[k] {
			t.Errorf("%v, was expected to be %v", result, expectedBoolSlice)
		}
	}
}

func TestGetColorsFromBytesReturnsCorrectForegroundColor(t *testing.T) {
	result, _ := getColorsFromBytes(0, 0, 0)
	if result.R != 0 || result.G != 0 || result.B != 0 {
		t.Errorf("RGB values were %d, %d, %d, expected %d, %d, %d", result.R, result.G, result.B, 0, 0, 0)
	}

	result, _ = getColorsFromBytes(255, 255, 255)
	if result.R != 255 || result.G != 255 || result.B != 255 {
		t.Errorf("RGB values were %d, %d, %d, expected %d, %d, %d", result.R, result.G, result.B, 255, 255, 255)
	}

	result, _ = getColorsFromBytes(40, 80, 120)
	if result.R != 40 || result.G != 80 || result.B != 120 {
		t.Errorf("RGB values were %d, %d, %d, expected %d, %d, %d", result.R, result.G, result.B, 40, 80, 120)
	}
}

func TestGetColorsFromBytesReturnsCorrectBackgroundColor(t *testing.T) {
	// On the verge, but still black
	_, result := getColorsFromBytes(168, 221, 55)
	if result.R != 0 || result.G != 0 || result.B != 0 {
		t.Errorf("RGB values were %d, %d, %d, expected %d, %d, %d", result.R, result.G, result.B, 0, 0, 0)
	}

	// On the verge, but still white
	_, result = getColorsFromBytes(168, 220, 55)
	if result.R != 255 || result.G != 255 || result.B != 255 {
		t.Errorf("RGB values were %d, %d, %d, expected %d, %d, %d", result.R, result.G, result.B, 255, 255, 255)
	}

	_, result = getColorsFromBytes(255, 255, 255)
	if result.R != 0 || result.G != 0 || result.B != 0 {
		t.Errorf("RGB values were %d, %d, %d, expected %d, %d, %d", result.R, result.G, result.B, 0, 0, 0)
	}

	_, result = getColorsFromBytes(0, 0, 0)
	if result.R != 255 || result.G != 255 || result.B != 255 {
		t.Errorf("RGB values were %d, %d, %d, expected %d, %d, %d", result.R, result.G, result.B, 255, 255, 255)
	}
}
