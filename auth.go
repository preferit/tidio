package tidio

import (
	"encoding/base64"
	"net/http"
)

func BasicAuth(r *http.Request, cred Credentials) (*http.Request, error) {
	if cred.account == "" { // anonymous
		return r, nil
	}
	plain := []byte(cred.account + ":" + cred.secret)
	v := base64.StdEncoding.EncodeToString(plain)
	r.Header.Set("Authorization", "Basic "+v)
	return r, nil
}
