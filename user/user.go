package user

import (
  "crypto/md5"
  "io"
  "fmt"
  "os"
  "github.com/duranmla/clirescue/cmdutil"
)

var (
  Stdout *os.File = os.Stdout
)

// the idea is that the same identicon is going to be retrieved using email, ip or token
type User struct {
  name          string `json:"name"`
  email         string `json:"email"`
  hash          string
}

func New(name, email string) *User {
	user := &User{name: name, email: email}
  user.hash = user._getStringMD5(email)
  return user
}

func RequestCredentials() (name, email string){
  fmt.Fprint(Stdout, "name: ")
  name = cmdutil.ReadLine()
  fmt.Fprint(Stdout, "email: ")
  email = cmdutil.ReadLine()

  return name, email
}

func (user *User) _getStringMD5(str string) string {
	hash := md5.New()
	io.WriteString(hash, str)

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (user *User) String() string  {
  return fmt.Sprintf("Hi %s!\nyour email is: %s\nwe've generated an image with: %s\n", user.name, user.email, user.hash)
}
