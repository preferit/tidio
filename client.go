package tidio

import (
	"fmt"
	"io"
	"net/http"
)

func NewClient(cred Credentials, options ...ClientOption) *Client {
	client := Client{
		Client: http.DefaultClient,
		cred:   cred,
	}
	for _, option := range options {
		option.ForClient(&client)
	}
	return &client
}

type Client struct {
	*http.Client
	cred Credentials
	host string

	err  error          // last error
	req  *http.Request  // last request
	resp *http.Response // last response
}

func (me *Client) Request() *http.Request   { return me.req }
func (me *Client) Response() *http.Response { return me.resp }

func (me *Client) CreateTimesheet(loc string, body io.Reader) error {
	if me.err != nil {
		return nil
	}
	me.newRequest("POST", me.host+loc, body)
	me.send()
	me.checkStatusCode(201)
	return me.err
}

// checkStatusCode
func (me *Client) checkStatusCode(exp int) {
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
func (me *Client) newRequest(method, path string, body io.Reader) {
	me.req, me.err = http.NewRequest(method, path, body)
}

func (me *Client) send() {
	if me.err != nil {
		return
	}
	if me.req == nil {
		panic("cannot send nil")
	}
	me.req.Header = me.cred.BasicAuth()
	me.resp, me.err = me.Do(me.req)
}
