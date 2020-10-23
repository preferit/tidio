package tidio

import (
	"io"
	"net/http"

	"github.com/gregoryv/ant"
)

// NewAPI returns the API for the given host.
//
// In the future some kind of validation might be put here that the
// host is compatible with the given implementation.
func NewAPI(host string, settings ...ant.Setting) *API {
	api := API{host: host}
	ant.MustConfigure(&api, settings...)
	return &api
}

// API provides http request builders for the tidio service
// The requests returned should be valid and complete.
type API struct {
	host string
	cred Credentials
}

// SetCredentials
func (me *API) SetCredentials(c Credentials) { me.cred = c }

/*

todo

- Should the API provide validation on this end? probably even
though it's needed on the receiving end.

-Where should we put the model parsing of the response?

one request could be parsed in multiple ways depending on what is
needed. If we document a full response a model for parsing it may not
be necessary, though could be provided for convenience.

*/

func (me API) CreateTimesheet(loc string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest("POST", me.url(loc), body)
	me.cred.BasicAuth(r.Header)
	return r, err
}

func (me API) ReadTimesheet(loc string) (*http.Request, error) {
	return http.NewRequest("GET", me.url(loc), nil)
}

func (me API) url(path string) string {
	return me.host + path
}
