package tidio

import (
	"crypto/sha512"
	"fmt"
)

func NewKey(secret, accountPath string) *Key {
	return &Key{
		Secret:      sha512.Sum512_256([]byte(secret)),
		AccountPath: accountPath,
	}
}

// Key links an account with a secret
type Key struct {
	Secret      [32]byte // SHA512_256 encrypted
	AccountPath string
}

// Check
func (me *Key) Check(secret string) error {
	if me.Secret != sha512.Sum512_256([]byte(secret)) {
		return fmt.Errorf("Check: incorrect secret")
	}
	return nil
}
