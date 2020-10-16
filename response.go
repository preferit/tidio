package tidio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gregoryv/rs"
)

// Response defines a http response ready to be written.
type Response struct {
	view       interface{}
	statusCode int
	sys        *rs.System
}

func (me *Response) authenticate(r *http.Request) (*rs.Account, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return rs.Anonymous, nil
	}

	name, secret, ok := r.BasicAuth()
	if !ok {
		return rs.Anonymous, fmt.Errorf("authentication failed")
	}
	asRoot := rs.Root.Use(me.sys)
	cmd := rs.NewCmd("/bin/secure", "-c", "-a", name, "-s", secret)
	if err := asRoot.Run(cmd); err != nil {
		return rs.Anonymous, err
	}
	var acc rs.Account
	err := asRoot.LoadAccount(&acc, name)
	return &acc, err
}

func (resp *Response) Build(r *http.Request) error {
	// todo perhaps add help as a resource
	if r.URL.Path == "/" {
		return resp.End(http.StatusOK, NewHelpView())
	}

	acc, err := resp.authenticate(r)
	if err != nil {
		return resp.Fail(http.StatusUnauthorized, err)
	}
	asAcc := acc.Use(resp.sys)

	// Check resource access permissions
	res, err := asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		return resp.Fail(http.StatusUnauthorized, err)
	}

	// Serve the specific method
	switch r.Method {
	case "GET":
		if res == nil {
			return resp.Fail(http.StatusNotFound, err)
		}
		if res.IsDir() == nil {
			cmd := rs.NewCmd("/bin/ls", "-json", "-json-name", "resources", r.URL.Path)
			var buf bytes.Buffer
			cmd.Out = &buf
			asAcc.Run(cmd)
			return resp.End(http.StatusOK, &buf)
		}
		res, err := asAcc.Open(r.URL.Path)
		if err != nil {
			return resp.Fail(http.StatusUnauthorized, err)
		}
		return resp.End(http.StatusOK, res)

	case "POST":
		if r.Body != nil {
			defer r.Body.Close()
		}
		res, err := asAcc.Create(r.URL.Path)
		if err != nil {
			// FIXME if write permission error use Forbidden
			return resp.Fail(http.StatusBadRequest, err)
		}
		resp.statusCode = http.StatusCreated
		io.Copy(res, r.Body)
		res.Close() // important to flush the data
		return resp.End(http.StatusCreated)
	default:
		return resp.Fail(http.StatusMethodNotAllowed, fmt.Errorf("Method not allowed"))
	}
}

// Fail sets the given status code on the response and returns the
// given error.
func (me *Response) Fail(code int, err error) error {
	me.statusCode = code
	return err
}

// End sets a 2xx code, optional view and returns nil.
func (me *Response) End(code int, view ...interface{}) error {
	me.statusCode = code
	if len(view) > 0 {
		me.view = view[0]
	}
	return nil
}

// WriteError
func (me *Response) WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "plain/text")
	if me.statusCode == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", `Basic realm="tidio"`)
	}
	w.WriteHeader(me.statusCode)
	w.Write([]byte(err.Error()))
}

// Send
func (me *Response) Send(w http.ResponseWriter) {
	w.WriteHeader(me.statusCode)
	switch view := me.view.(type) {
	case io.ReadCloser:
		io.Copy(w, view)
		view.Close()
		return
	case io.WriterTo:
		view.WriteTo(w)
		return
	}
}
