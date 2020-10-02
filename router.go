package tidio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/rs"
)

func NewRouter(sys *rs.System) *Router {
	nop := fox.NewSyncLog(ioutil.Discard)

	return &Router{
		Logger: nop,
		warn:   nop.Log,
		sys:    sys,
	}
}

type Router struct {
	fox.Logger
	warn func(...interface{})

	sys *rs.System
}

func (me *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp Response
	err := me.endpoint(&resp, r)
	if err != nil {
		// Handle any errors
		w.Header().Set("Content-Type", "plain/text")
		if resp.statusCode == http.StatusUnauthorized {
			w.Header().Set("WWW-Authenticate", `Basic realm="tidio"`)
		}
		w.WriteHeader(resp.statusCode)
		w.Write([]byte(err.Error()))

		return
	}
	// Render ok response
	w.WriteHeader(resp.statusCode)
	switch view := resp.view.(type) {
	case io.ReadCloser:
		io.Copy(w, view)
		view.Close()
		return
	case io.WriterTo:
		view.WriteTo(w)
		return
	}
}

func (me *Router) endpoint(resp *Response, r *http.Request) error {

	// TODO perhaps add help as a resource
	if r.URL.Path == "/" {
		return resp.End(http.StatusOK, NewHelpView())
	}

	// TODO design routing as a chained responsibility
	//
	// Select response format by accept header
	// Authorize
	// Check resources
	// Find command for Method+Mimetype
	// Exec and respond with
	acc, err := me.authenticate(r)
	if err != nil {
		return resp.Fail(http.StatusUnauthorized, err)
	}
	asAcc := acc.Use(me.sys)

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

func (me *Router) authenticate(r *http.Request) (*rs.Account, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return rs.Anonymous, nil
	}

	name, secret, ok := r.BasicAuth()
	me.Log("authenticate ", name)
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
