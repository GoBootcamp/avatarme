package avatar

import (
  "fmt"
  "github.com/duranmla/avatarme/user"
)

type Avatar struct {
  *user.User
  ink     string
  pixels  [12][2]int // [[0,1], [4,5]...]
}

func New(name, email string) *Avatar {
  avatar := &Avatar{User: user.New(name, email)}
  avatar._getAvatarColor()
  avatar._getAvatarPixels()

  fmt.Println(avatar)
  return avatar
}

func (avatar *Avatar) _getAvatarColor() {
  avatar.ink = avatar.hash[26:] // last 6 characteres (Hexcolor)
}

func (avatar *Avatar) _getAvatarPixels() {
  source := avatar.hash[:26] // get first 26 characteres

  for i:=0; i<=(len(source)-2); i+2 {
    if i%2 == 0 {
      avatar.pixels[i/2] = [2]int{avatar.hash[i:(i+1)], avatar.hash[(i+1):(i+2)]}
    }
  }
}

func (avatar *Avatar) String() string  {
  return fmt.Sprintf("Hi %s!\nyour email is: %s\nwe've generated an image with: %s\n\n color assigned was: %s\n", avatar.name, avatar.email, avatar.hash, avatar.ink)
}
