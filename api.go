package tidio

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/fox"
)

// NewAPI returns the API for the given host.
//
// In the future some kind of validation might be put here that the
// host is compatible with the given implementation.
func NewAPI(host string, settings ...ant.Setting) *API {
	api := API{
		Logger: fox.NewSyncLog(ioutil.Discard),
		host:   host,
		client: http.DefaultClient,
		auth:   BasicAuth,
	}
	ant.MustConfigure(&api, settings...)
	return &api
}

// API provides http request builders for the tidio service
// The requests returned should be valid and complete.
type API struct {
	fox.Logger
	host   string
	client *http.Client
	asJohn Credentials
	auth   func(*http.Request, Credentials) (*http.Request, error)

	// last api
	Request *http.Request
	Err     error
}

func (me *API) SetLogger(l fox.Logger) { me.Logger = l }

// Auth applies credentials to the request and sets them as last
// values on the api.
func (me *API) Auth(r *http.Request, err error) {
	me.Request = r
	me.Err = err
	if err != nil {
		return
	}
	me.Request, me.Err = me.auth(r, me.asJohn)
}

func (me *API) CreateTimesheet(loc string, body io.Reader) *API {
	me.Auth(http.NewRequest("POST", me.url(loc), body))
	return me
}

func (me *API) ReadTimesheet(loc string) *API {
	me.Auth(http.NewRequest("GET", me.url(loc), nil))
	return me
}

func (me *API) SetCredentials(c Credentials) { me.asJohn = c }

func (me *API) url(path string) string {
	return me.host + path
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
	resp, err := me.client.Do(r)
	if err != nil {
		me.Log(r.Method, r.URL, err)
		return resp, err
	}
	me.Log(r.Method, r.URL, resp.StatusCode)
	return resp, nil
}

/*

todo

- Should the API provide validation on this end? probably even
though it's needed on the receiving end.

-Where should we put the model parsing of the response?

one request could be parsed in multiple ways depending on what is
needed. If we document a full response a model for parsing it may not
be necessary, though could be provided for convenience.

*/
