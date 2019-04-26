package avatarme

import (
	"testing"
	"io/ioutil"
)

func TestUniqueHash(t *testing.T) {
  g, _ := UniqueHash([]byte("test msg"))
  if g != 16644892400343469032 {
    t.Errorf("UniqueHash([]byte(\"test msg\")) = %d, want 16644892400343469032", g)
  }
	
  f, _ := ioutil.ReadFile("testfile.txt")
  g, _ = UniqueHash(f)
  if g != 16644892400343469032 {
    t.Errorf("UniqueHash([]byte(\"test msg\")) = %d, want 16644892400343469032", g)
  }

  g, _ = UniqueHash([]byte(""))
  if g != 14695981039346656037 {
    t.Errorf("UniqueHash([]byte(\"\")) = %d, want 14695981039346656037", g)
  }

  g, _ = UniqueHash([]byte{0})
  if g != 12638153115695167455 {
    t.Errorf("UniqueHash([]byte{0}) = %d, want 12638153115695167455", g)
  }

  g, err := UniqueHash([]byte{})
  if err == nil {
    t.Errorf("UniqueHash([]byte{}) = %d, want Error about empty array, g)
  }
}