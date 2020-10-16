package tidio

import (
	"io"
	"net/http"
)

func NewAPI(host string) *API {
	return &API{host: host}
}

type API struct {
	host string
}

func (me API) CreateTimesheet(loc string, body io.Reader) (*http.Request, error) {
	return http.NewRequest("POST", me.url(loc), body)
}

func (me API) ReadTimesheet(loc string) (*http.Request, error) {
	return http.NewRequest("GET", me.url(loc), nil)
}

func (me API) url(path string) string {
	return me.host + path
}
