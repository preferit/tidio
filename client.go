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
}

func (me *Client) CreateTimesheet(loc string, body io.Reader) error {
	r, err := http.NewRequest("POST", me.host+loc, body)
	if err != nil {
		me.Log(err)
		return err
	}
	r.Header = me.cred.BasicAuth()
	resp, err := me.Do(r)
	if err != nil {
		me.Log(r.Method, r.URL, err)
		return err
	}
	me.Log(r.Method, r.URL, resp.StatusCode)
	return checkStatusCode(resp, 201)
}

// checkStatusCode
func checkStatusCode(resp *http.Response, exp int) error {
	if resp.StatusCode != exp {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return nil
}
