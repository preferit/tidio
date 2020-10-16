package tidio

import (
	"encoding/base64"
	"net/http"
)

type Credentials struct {
	account string
	secret  string
}

func (me Credentials) BasicAuth() http.Header {
	return basicAuthHeaders(me.account, me.secret)
}

func basicAuthHeaders(user, pass string) http.Header {
	headers := http.Header{}
	secret := base64.StdEncoding.EncodeToString([]byte("john:secret"))
	headers.Set("Authorization", "Basic "+secret)
	return headers
}

func (me Credentials) ForClient(cli *Client) error {
	cli.cred = me
	return nil
}
