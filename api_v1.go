package tidio

import (
	"fmt"
	"io"
	"net/http"
)

type APIv1 interface {
	CreateTimesheet(path string, timesheet io.Reader) error
}

func NewAPIv1(host string, cred Credentials) APIv1 {
	return &v1{
		Client: http.DefaultClient,
		cred:   cred,
		host:   host,
	}
}

type v1 struct {
	*http.Client
	cred Credentials
	host string

	err  error          // last error
	req  *http.Request  // last request
	resp *http.Response // last response
}

func (me *v1) Request() *http.Request   { return me.req }
func (me *v1) Response() *http.Response { return me.resp }

func (me *v1) CreateTimesheet(loc string, body io.Reader) error {
	if me.err != nil {
		return nil
	}
	me.newRequest("POST", me.host+loc, body)
	me.send()
	me.checkStatusCode(201)
	return me.err
}

// checkStatusCode
func (me *v1) checkStatusCode(exp int) {
	if me.err != nil {
		return
	}
	if me.resp == nil {
		panic("response is nil")
	}
	if me.resp.StatusCode != exp {
		me.err = fmt.Errorf("unexpected status code: %v", me.resp.StatusCode)
	}
}

// newRequest
func (me *v1) newRequest(method, path string, body io.Reader) {
	me.req, me.err = http.NewRequest(method, path, body)
}

func (me *v1) send() {
	if me.err != nil {
		return
	}
	if me.req == nil {
		panic("cannot send nil")
	}
	me.req.Header = me.cred.BasicAuth()
	me.resp, me.err = me.Do(me.req)
}
