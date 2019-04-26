package avatarme

import "hash/fnv"

// uniqueHash returns unsigned 64bit int and uses fnv hashing algorithm.
func uniqueHash(b []byte) (uint64, error) {
  if len(b) == 0 {
    return uint64(0), errors.New("tried to initalise empty byte array")
  }
  
  newHash := fnv.New64a()
  _, err := newHash.Write(b)
  if err != nil {
    return uint64(0), err
  }
  
  return newHash.Sum64(), nil
}
