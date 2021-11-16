package tidio

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gregoryv/ant"
	"github.com/gregoryv/rs"
)

func (me *System) Router() *mux.Router {
	r := mux.NewRouter()
	r.Methods("GET").HandlerFunc(me.serveRead)
	r.Methods("POST").HandlerFunc(me.serveCreate)
	log := Register(r)
	r.Use(NewRouteLog(log).Middleware)
	if Conf.Debug() {
		r.Use(NewTracer(log).Middleware)
	}
	return r
}

func (me *System) serveRead(w http.ResponseWriter, r *http.Request) {
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

func (me *System) serveCreate(w http.ResponseWriter, r *http.Request) {
	acc, err := authenticate(me.sys, r)
	log := Log(r)
	log.Info(acc.Name)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		log.Info(err)
		return
	}
	asAcc := acc.Use(me.sys)
	asAcc.SetAuditer(log)

	log.Info(acc)
	// Check resource access permissions
	_, err = asAcc.Stat(r.URL.Path)
	if acc == rs.Anonymous && err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, err)
		log.Error(err)
		return
	}

	// Serve the specific method
	defer r.Body.Close()
	res, err := asAcc.Create(r.URL.Path)
	if err != nil {
		log.Error(err)
		// todo probably return an error here
	}
	// TODO when sharing is implemented and accounts have read but not write permissions
	// Create will fail

	io.Copy(res, r.Body)
	res.Close() // important to flush the data
	w.WriteHeader(http.StatusCreated)
	go me.PersistState()
}

func NewBasicAuth(c *Credentials) *BasicAuth {
	return &BasicAuth{cred: c}
}

type BasicAuth struct {
	cred *Credentials
}

func (me *BasicAuth) Set(v interface{}) error {
	if me.cred == nil { // anonymous
		return nil
	}
	switch v := v.(type) {
	case *http.Request:
		plain := []byte(me.cred.account + ":" + me.cred.secret)
		b := base64.StdEncoding.EncodeToString(plain)
		v.Header.Set("Authorization", "Basic "+b)
		return nil
	default:
		return ant.SetFailed(v, me)
	}
}

func authenticate(sys *rs.System, r *http.Request) (*rs.Account, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return rs.Anonymous, nil
	}

	name, secret, ok := r.BasicAuth()
	if !ok {
		return rs.Anonymous, fmt.Errorf("authentication failed")
	}

	asRoot := rs.Root.Use(sys)
	asRoot.SetAuditer(Log(r))
	cmd := rs.NewCmd("/bin/secure", "-c", "-a", name, "-s", secret)
	if err := asRoot.Run(cmd); err != nil {
		return rs.Anonymous, err
	}
	var acc rs.Account
	err := asRoot.LoadAccount(&acc, name)
	return &acc, err
}

func NewTracer(dst *LogPrinter) *Tracer {
	return &Tracer{dst: dst}
}

type Tracer struct {
	dst *LogPrinter
}

func (me *Tracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := Register(r).Buf()
		log.lgr.SetPrefix("TRACE: ")
		next.ServeHTTP(w, r)
		if log.Failed() {
			me.dst.Info(log.FlushString())
		}
		Unregister(r)
	})
}

// NewRouteLog returns RouteLog using the package logger.
func NewRouteLog(dst *LogPrinter) *RouteLog {
	return &RouteLog{LogPrinter: dst}
}

type RouteLog struct {
	*LogPrinter
}

// Middleware returns a middleware that logs request and response status.
func (me *RouteLog) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{w, 200}
		next.ServeHTTP(&rec, r)
		me.Info(r.Method, r.URL, rec.status, time.Since(start))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
	rec.status = code
}
