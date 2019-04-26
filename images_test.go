package avatarme

import (
  "testing"
  "image/color"
)

func TestDeduceShape(t *testing.T) {
  var temp [5][5]int
  if a := deduceShape(uint64(1)); a != temp {
	t.Errorf("expected zeroed array")
  }
}

func TestGetColors(t *testing.T) {
  fg, bg := getColors(uint64(1)) // should take first 24 bits
  if x := fmt.Sprintf("%d%d%d", fg.R, fg.G, fg.B); x != "100" {
    t.Errorf("expected value 100, not %d%d%d", fg.R, fg.G, fg.B)
  }
  if x := fmt.Sprintf("%d%d%d", bg.R, bg.G, bg.B); x != "000" {
    t.Errorf("expected value 000, not %d%d%d", bg.R, bg.G, bg.B)
  }

  fg, bg = getColors(uint64(16777216)) // should take next 24 bits
  if x := fmt.Sprintf("%d%d%d", fg.R, fg.G, fg.B); x != "000" {
    t.Errorf("expected value 100, not %d%d%d", fg.R, fg.G, fg.B)
  }
  if x := fmt.Sprintf("%d%d%d", bg.R, bg.G, bg.B); x != "100" {
    t.Errorf("expected value 000, not %d%d%d", bg.R, bg.G, bg.B)
  }
}

func TestDeduceRGB(t *testing.T) {
  r, g, b := deduceRGB(uint64(0))
  if r != 0 && g != 0 && b != 0 {
    t.Errorf("expected 0 0 0 in deduceRGB instead of %d %d %d", r, g, b)
  }

  r, g, b = deduceRGB(uint64(255)) // 00000000 00000000 11111111
  if r != 255 && g != 0 && b != 0 {
    t.Errorf("expected 255 0 0 in deduceRGB instead of %d %d %d", r, g, b)
  }

  r, g, b = deduceRGB(uint64(65280)) // 00000000 11111111 00000000
  if r != 0 && g != 255 && b != 0 {
    t.Errorf("expected 0 255 0 in deduceRGB instead of %d %d %d", r, g, b)
  }

  r, g, b = deduceRGB(uint64(16711680)) // 11111111 00000000 00000000
  if r != 0 && g != 0 && b != 255 {
    t.Errorf("expected 0 0 255 in deduceRGB instead of %d %d %d", r, g, b)
  }
}