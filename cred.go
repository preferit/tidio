package tidio

import (
	"encoding/base64"
	"net/http"
)

type Credentials struct {
	account string
	secret  string
}

func (me Credentials) BasicAuth(h *http.Header) {
	plain := []byte(me.account + ":" + me.secret)
	v := base64.StdEncoding.EncodeToString(plain)
	h.Set("Authorization", "Basic "+v)
}
