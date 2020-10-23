package tidio

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gregoryv/ant"
)

// NewAPI returns the API for the given host.
//
// In the future some kind of validation might be put here that the
// host is compatible with the given implementation.
func NewAPI(host string, settings ...ant.Setting) *API {
	api := API{
		host:   host,
		client: http.DefaultClient,
	}
	api.SetCredentials(nil)
	ant.MustConfigure(&api, settings...)
	return &api
}

// API provides http request builders for the tidio service
// The requests returned should be valid and complete.
type API struct {
	Logger
	host   string
	client *http.Client
	auth   ant.Setting // applied

	// last api
	Request *http.Request
}

func (me *API) CreateTimesheet(path string, body io.Reader) *API {
	r := me.newRequest("POST", path, body)
	me.Auth(r)
	return me
}

func (me *API) ReadTimesheet(path string) *API {
	r := me.newRequest("GET", path, nil)
	me.Auth(r)
	return me
}

func (me *API) SetCredentials(c *Credentials) {
	me.auth = NewBasicAuth(c)
}

// Auth applies credentials to the request and sets them as last
// values on the api.
func (me *API) Auth(r *http.Request) {
	if r == nil {
		return
	}
	err := ant.Configure(r, me.auth)
	me.warn(err)
}

// newRequest
func (me *API) newRequest(method, path string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, me.host+path, body)
	me.warn(err)
	me.Request = r
	return r
}

// MustSend
func (me *API) MustSend() *http.Response {
	r, err := me.Send()
	if err != nil {
		panic(err)
	}
	return r
}

func (me *API) Send() (*http.Response, error) {
	r := me.Request
	if r == nil {
		return nil, fmt.Errorf("missing API request")
	}
	resp, err := me.client.Do(r)
	if err != nil {
		me.Log(r.Method, r.URL, err)
		return resp, err
	}
	me.Log(r.Method, r.URL, resp.StatusCode)
	return resp, nil
}

// warn logs non nil errors
func (me *API) warn(err error) {
	if err == nil {
		return
	}
	me.Log(err)
}
