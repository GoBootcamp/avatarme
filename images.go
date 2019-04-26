package avatarme

import (
  "bytes"
  "image"
  "image/color"
  "image/png"
  "bytes"
)

// genImg uses seed to generate 5x5 png image that provides symetric Identicon
func genImg(u uint64) []byte {
  height, width := 5, 5
  img := image.NewRGBA(image.Rect(0, 0, height, width))
  
  shape := deduceShape(u)
  fgColor, bgColor := getColors(u)
  
  var currentColor color.RGBA
  for y,x := range shape {
    for f, g := range x {
      currentColor = fgColor
      if g == 0 {
        currentColor = bgColor
      }
      
      img.Set(y, f, currentColor)
    }
  }
  
  var buff bytes.Buffer
  png.Encode(&buff, img)
  return buff.Bytes()
}

// deduceShape uses all 64bits to toggle 15 values between 0 and 1.
// After that it generates second symetrical part of the shape. 
func deduceShape(u uint64) [5][5]int {
  var basicShape [5][5]int
  var x,y int
  // it draws only first three columns on x axis because
  // column 4th and 5th will be cloned from column 1st and 2nd
  for i := uint64(0); i < 64; i++ {
    basicShape[x][y] = 0
    if (u & (1 << i)) != 0 {
      basicShape[x][y] = 1
    }
    
    y++
    if y > 3 { x++; y = 0 }
    if x > 3 { x = 0 }
  }
  
  for y, x := range basicShape[0] {
    basicShape[4][y] = x
  }
  
  for y, x := range basicShape[1] {
    basicShape[3][y] = x
  }
  
  return basicShape
}

// getColors returns two RGB sets, foreground color uses first 24bits while background color uses next 24bits.
func getColors(u uint64) (color.RGBA, color.RGBA) {
  fg := (0xFFFFFF & u)
  bg := (0xFFFFFF000000 & u) >> 24
  redFg, greenFg, blueFg := deduceRGB(fg)
  redBg, greenBg, blueBg := deduceRGB(bg)
  return color.RGBA{redFg, greenFg, blueFg, 0xFF}, color.RGBA{redBg, greenBg, blueBg, 0xFF}
}

// deduceRGB splits 24bit input between 3 values.
func deduceRGB(u uint64) (uint8, uint8, uint8) {
  return uint8(u & 255), uint8((u >> 8) & 255), uint8((u >> 16) & 255)
}