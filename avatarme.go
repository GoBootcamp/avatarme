// Package avatarme provides an easy way to create simple Identicons.
package avatarme

import (
  "encoding/base64"
  "os"
  "strings"
  "image"
  "image/png"
)

// Identicon stores the image encoded to base64 and filename of a file it should be saved to.
type Identicon struct {
  base64    string
  filename  string
}

// New creates a new Identicon object that accepts input data and filename of future file as arguments.
func New(inp []byte, f filename) *Identicon {
  hashValue := uniqueHash(inp)
  identicon := genImg(hashValue)
  return &Identicon{base64: base64.StdEncoding.EncodeToString(identicon), filename: f}
}

// Base64 returns png image as base64 encoded string
func (i *Identicon) Base64() string {
  return i.base64
}

// Draw creates an 5x5 Image containing an identicon that is saved as a png to the location specified in Identicon object. 
func (i *Identicon) Draw() error {
  reader := base64.NewDecorder(base64.StdEncoding, strings.NewReader(i.base64()))
  img, _, err := image.Decode(reader)
  if err != nil {
    return err
  }
  
  f, err := os.OpenFile(i.filename)
  if err != nil {
    return err
  }
  defer f.Close()
  
  err = png.Encode(f, img)
  if err != nil {
    return err
  }
  
  return nil
}
