package identicon

import (
	"crypto/sha256"
	"encoding/hex"
)

type Identicon struct {
	identifier, hashString string
	hashBytes [32]byte
}

func New(identifier string) *Identicon {
	hash := sha256.Sum256([]byte(identifier))
	hashString := hash[:]
	identicon := &Identicon{
		identifier: identifier,
		hashBytes: hash,
		hashString: hex.EncodeToString(hashString),
	}
	return identicon
}