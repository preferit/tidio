package tidio

import (
	"encoding/base64"
	"net/http"
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

func (me Credentials) BasicAuth(h *http.Header) {
	if me.account == "" { // anonymous
		return
	}
	plain := []byte(me.account + ":" + me.secret)
	v := base64.StdEncoding.EncodeToString(plain)
	h.Set("Authorization", "Basic "+v)
}
