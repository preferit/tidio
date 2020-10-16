package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gregoryv/fox"
)

func NewClient(settings ...Setting) *Client {
	client := Client{
		Client: http.DefaultClient,
		Logger: fox.NewSyncLog(ioutil.Discard),
	}
	for _, setting := range settings {
		err := setting.Set(&client)
		if err != nil {
			panic(err) // or client.err = err
		}
	}
	return &client
}

type Client struct {
	*http.Client
	fox.Logger
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
		return me.err
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
	me.Log(me.req.Method, me.req.URL, me.resp.StatusCode)
}
