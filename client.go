package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gregoryv/fox"
)

func NewClient(settings ...Setting) *Client {
	client := Client{
		Client: http.DefaultClient,
		Logger: fox.NewSyncLog(ioutil.Discard),
	}
	err := client.Use(settings...)
	if err != nil {
		panic(err)
	}
	return &client
}

type Client struct {
	*http.Client
	fox.Logger
	check func(...interface{})

	cred *Credentials
	host string
}

// Use configures client with new settings.
func (me *Client) Use(settings ...Setting) error {
	for _, setting := range settings {
		err := setting.Set(me)
		if err != nil {
			me.handle(&err)
			return err
		}
	}
	return nil
}

// handle
func (me *Client) handle(err *error) error {
	if *err != nil && me.check != nil {
		me.check(*err)
	}
	return *err
}

// send
func (me *Client) Send(r *http.Request, cred *Credentials) (*http.Response, error) {
	if cred != nil {
		cred.BasicAuth(&r.Header)
	}
	fullURL, _ := url.Parse(me.host + r.URL.String())
	r.URL = fullURL
	resp, err := me.Do(r)
	if err != nil {
		me.Log(r.Method, r.URL, err)
		return resp, err
	}
	me.Log(r.Method, r.URL, resp.StatusCode)
	return resp, nil
}

// todo client should send requests with authorization and normalize
// error handling of responses
func (me *Client) CreateTimesheet(loc string, body io.Reader) (err error) {
	defer me.handle(&err)
	api := API{}

	r, err := api.CreateTimesheet(loc, body)
	if err != nil {
		return
	}
	resp, err := me.Send(r, me.cred)
	if err != nil {
		return
	}
	return checkStatusCode(resp, 201)
}

func (me *Client) ReadTimesheet(loc string) (body io.ReadCloser, err error) {
	defer me.handle(&err)
	r, err := http.NewRequest("GET", me.host+loc, nil)
	if err != nil {
		me.Log(err)
		return
	}
	if me.cred != nil {
		me.cred.BasicAuth(&r.Header)
	}
	resp, err := me.Do(r)
	if err != nil {
		me.Log(r.Method, r.URL, err)
		return
	}
	me.Log(r.Method, r.URL, resp.StatusCode)
	body = resp.Body
	err = checkStatusCode(resp, 200)
	return
}

// checkStatusCode
func checkStatusCode(resp *http.Response, exp int) error {
	if resp.StatusCode != exp {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return nil
}
