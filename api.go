package tidio

import (
	"io"
	"net/http"
)

type API struct {
	host string
	cred *Credentials
}

// CreateTimesheet
func (me API) CreateTimesheet(loc string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest("POST", me.host+loc, body)
	if err != nil {
		return nil, err
	}
	if me.cred != nil {
		me.cred.BasicAuth(&r.Header)
	}
	return r, nil
}
