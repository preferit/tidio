package tidio

import (
	"github.com/gregoryv/ant"
)

// Credentials provides ways to authenticate a requests via header
// manipulation. Zero value credentials means anonymous.
func NewCredentials(account, secret string) *Credentials {
	return &Credentials{
		account: account,
		secret:  secret,
	}
}

type Credentials struct {
	account string
	secret  string
}

func (me *Credentials) Set(v interface{}) error {
	switch v := v.(type) {
	case usesCredentials:
		v.SetCredentials(me)
	case *Service:
		v.AddAccount(me.account, me.secret)
	default:
		return ant.SetFailed(v, me)
	}
	return nil
}

type usesCredentials interface {
	SetCredentials(*Credentials)
}
