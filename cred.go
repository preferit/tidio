package tidio

import (
	"encoding/base64"
	"net/http"

	"github.com/gregoryv/ant"
)

// Credentials provides ways to authenticate a requests via header
// manipulation. Zero value credentials means anonymous.
func NewCredentials(account, secret string) Credentials {
	return Credentials{
		account: account,
		secret:  secret,
	}
}

type Credentials struct {
	account string
	secret  string
}

// todo decouple use of credentials from authentication method
// maybe add type AuthMethod interface {
func (me Credentials) BasicAuth(h http.Header) {
	if me.account == "" { // anonymous
		return
	}
	plain := []byte(me.account + ":" + me.secret)
	v := base64.StdEncoding.EncodeToString(plain)
	h.Set("Authorization", "Basic "+v)
}

func (me Credentials) Set(v interface{}) error {
	switch v := v.(type) {
	case usesCredentials:
		v.SetCredentials(me)
	default:
		return ant.SetFailed(v, me)
	}
	return nil
}

type usesCredentials interface {
	SetCredentials(Credentials)
}
