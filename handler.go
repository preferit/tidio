package tidio

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gregoryv/rs"
)

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	acc, err := me.authenticate(r)
	if err != nil {
		me.warn(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	asAcc := acc.Use(me.sys)

	// Check resource access permissions
	res, err := asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		textErr(w, http.StatusUnauthorized, err)
		return
	}

	// Serve the specific method
	switch r.Method {
	case "GET":
		if res.IsDir() == nil {
			cmd := rs.NewCmd("/bin/ls", "-json", "-json-name", "resources", r.URL.Path)
			cmd.Out = w
			asAcc.Run(cmd)
			return
		}
		res, err := asAcc.Open(r.URL.Path)
		if err != nil {
			me.warn(err)
			textErr(w, http.StatusUnauthorized, err)
			return
		}
		io.Copy(w, res)
	case "POST":
		if r.Body != nil {
			defer r.Body.Close()
		}
		res, err := asAcc.Create(r.URL.Path)
		if err != nil {
			textErr(w, http.StatusBadRequest, err)
			return
		}
		io.Copy(res, r.Body)
		res.Close() // important to flush the data
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func textErr(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte(err.Error()))
}

func (me *Service) authenticate(r *http.Request) (*rs.Account, error) {
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
