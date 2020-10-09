package tidio

import (
	"encoding/base64"
	"net/http"
)

type Credentials struct {
	account string
	secret  string
}

// Method
func (me Credentials) BasicAuth() http.Header {
	return basicAuthHeaders(me.account, me.secret)
}

func basicAuthHeaders(user, pass string) http.Header {
	headers := http.Header{}
	secret := base64.StdEncoding.EncodeToString([]byte("john:secret"))
	headers.Set("Authorization", "Basic "+secret)
	return headers
}
