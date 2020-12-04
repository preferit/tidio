package tidio

import (
	"bytes"
	"fmt"
	"io"
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
	trace.Log(acc.Name)
	trace.Log(err)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		return
	}
	asAcc := acc.Use(me.sys)
	asAcc.SetAuditer(trace)

	trace.Log(acc)
	// Check resource access permissions
	_, err = asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		return
	}

	// Serve the specific method

	defer r.Body.Close()
	res, err := asAcc.Create(r.URL.Path)
	trace.Log(err)
	// TODO when sharing is implemented and accounts have read but not write permissions
	// Create will fail

	io.Copy(res, r.Body)
	res.Close() // important to flush the data
	w.WriteHeader(http.StatusCreated)
	go me.PersistState()
}
