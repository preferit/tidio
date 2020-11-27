package tidio

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gregoryv/fox/foxhttp"
	"github.com/gregoryv/rs"
)

func (me *Service) Router() *mux.Router {
	r := mux.NewRouter()
	r.Methods("GET").HandlerFunc(me.serveRead)
	r.Methods("POST").HandlerFunc(me.serveCreate)
	r.Use(foxhttp.NewRouteLog(me).Middleware)
	return r
}

func (me *Service) serveRead(w http.ResponseWriter, r *http.Request) {
	trace, cleanup := newTrace(me, r)
	defer cleanup()

	if r.URL.Path == "/" {
		w.WriteHeader(200)
		NewHelpView().WriteTo(w)
		return
	}

	acc, err := authenticate(me.sys, r, trace)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		return
	}
	asAcc := acc.Use(me.sys)

	// Check resource access permissions
	res, err := asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		return
	}

	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
		return
	}
	if res.IsDir() == nil {
		cmd := rs.NewCmd("/bin/ls", "-json", "-json-name", "resources", r.URL.Path)
		var buf bytes.Buffer
		cmd.Out = &buf
		asAcc.Run(cmd)
		w.WriteHeader(200)
		io.Copy(w, &buf)
		return
	}
	// todo check read error here
	resource, _ := asAcc.Open(r.URL.Path)
	w.WriteHeader(200)
	io.Copy(w, resource)
}

func (me *Service) serveCreate(w http.ResponseWriter, r *http.Request) {
	trace, cleanup := newTrace(me, r)
	defer cleanup()

	acc, err := authenticate(me.sys, r, trace)
	trace.Println(acc.Name)
	trace.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		return
	}
	asAcc := acc.Use(me.sys)

	trace.Println(acc)
	// Check resource access permissions
	_, err = asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		//me.Log(err)
		fmt.Fprintln(w, err)
		return
	}

	// Serve the specific method

	defer r.Body.Close()
	res, err := asAcc.Create(r.URL.Path)
	trace.Println(err)
	// TODO when sharing is implemented and accounts have read but not write permissions
	// Create will fail

	io.Copy(res, r.Body)
	res.Close() // important to flush the data
	w.WriteHeader(http.StatusCreated)
}

func authenticate(sys *rs.System, r *http.Request, trace *log.Logger) (*rs.Account, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return rs.Anonymous, nil
	}
	trace.Println(h)

	name, secret, ok := r.BasicAuth()
	if !ok {
		return rs.Anonymous, fmt.Errorf("authentication failed")
	}

	asRoot := rs.Root.Use(sys)
	cmd := rs.NewCmd("/bin/secure", "-c", "-a", name, "-s", secret)
	trace.Println(cmd)
	if err := asRoot.Run(cmd); err != nil {
		return rs.Anonymous, err
	}
	var acc rs.Account
	err := asRoot.LoadAccount(&acc, name)
	return &acc, err
}
