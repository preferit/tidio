package tidio

import (
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

	cred *Credentials
	host string
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

// Send
func (me *Client) Send(r *http.Request, cred *Credentials) (*http.Response, error) {
	if cred != nil {
		cred.BasicAuth(&r.Header)
	}
	resp, err := me.Do(r)
	if err != nil {
		me.Log(r.Method, r.URL, err)
		return resp, err
	}
	me.Log(r.Method, r.URL, resp.StatusCode)
	return resp, nil
}
