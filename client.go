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
	cred  *Credentials
	host  string
}

// Use configures client with new settings.
func (me *Client) Use(settings ...Setting) error {
	for _, setting := range settings {
		err := setting.Set(me)
		if err != nil {
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

func (me *Client) CreateTimesheet(loc string, body io.Reader) (err error) {
	defer me.handle(&err)
	r, err := http.NewRequest("POST", me.host+loc, body)
	if err != nil {
		me.Log(err)
		return
	}
	if me.cred != nil {
		r.Header = me.cred.BasicAuth()
	}
	resp, err := me.Do(r)
	if err != nil {
		me.Log(r.Method, r.URL, err)
		return
	}
	me.Log(r.Method, r.URL, resp.StatusCode)
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
		r.Header = me.cred.BasicAuth()
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
