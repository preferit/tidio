package tidio

import (
	"io/ioutil"
	"net/http"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/fox"
)

func NewClient(settings ...ant.Setting) *Client {
	client := Client{
		Client: http.DefaultClient,
		Logger: fox.NewSyncLog(ioutil.Discard),
	}
	ant.MustConfigure(&client, settings...)
	return &client
}

type Client struct {
	*http.Client
	fox.Logger
	check func(...interface{})
}

func (me *Client) SetLogger(l fox.Logger) { me.Logger = l }

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
