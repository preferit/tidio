package tidio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gregoryv/fox/foxhttp"
	"github.com/gregoryv/rs"
)

func (me *Service) Router() *mux.Router {
	r := mux.NewRouter()
	r.Methods("GET").HandlerFunc(me.serveRead)
	r.Methods("POST").HandlerFunc(me.serveCreate)
	log := RLog(r)
	log.SetOutput(ioutil.Discard)
	r.Use(
		foxhttp.NewRouteLog(log).Middleware,
		NewTracer(log).Middleware,
	)
	return r
}

func NewTracer(dst *LogPrinter) *tracer {
	return &tracer{dst: dst}
}

type tracer struct {
	dst *LogPrinter
}

func (me *tracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := RLog(r).Buf()
		log.lgr.SetPrefix("TRACE: ")
		next.ServeHTTP(w, r)
		if log.Failed() {
			me.dst.Info(log.FlushString())
		}
		Unreg(r)
	})
}

func (me *Service) serveRead(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		w.WriteHeader(200)
		NewHelpView().WriteTo(w)
		return
	}

	acc, err := authenticate(me.sys, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		Log(r).Error(err)
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
	trace := RLog(r).Buf()
	defer Unreg(r)
	acc, err := authenticate(me.sys, r)
	trace.Info(acc.Name)
	trace.Info(err)
	defer func() {
		if err != nil {
			Log(me).Info("trace\n", trace.FlushString())
		}
	}()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		return
	}
	asAcc := acc.Use(me.sys)
	asAcc.SetAuditer(trace)

	trace.Info(acc)
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
	trace.Info(err)
	// TODO when sharing is implemented and accounts have read but not write permissions
	// Create will fail

	io.Copy(res, r.Body)
	res.Close() // important to flush the data
	w.WriteHeader(http.StatusCreated)
	go me.PersistState()
}
