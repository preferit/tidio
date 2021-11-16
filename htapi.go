package tidio

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/apidoc"
	"github.com/gregoryv/web/toc"

	"github.com/gorilla/mux"
	"github.com/gregoryv/ant"
	"github.com/gregoryv/rs"
)

func NewRouter(sys *System) *mux.Router {
	r := mux.NewRouter()
	ht := HTAPI{System: sys}
	r.Methods("GET").HandlerFunc(ht.serveRead)
	r.Methods("POST").HandlerFunc(ht.serveCreate)
	log := Register(r)
	r.Use(NewRouteLog(log).Middleware)
	if Conf.Debug() {
		r.Use(NewTracer(log).Middleware)
	}
	return r
}

type HTAPI struct {
	*System
}

func (me *HTAPI) serveRead(w http.ResponseWriter, r *http.Request) {
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

func (me *HTAPI) serveCreate(w http.ResponseWriter, r *http.Request) {
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

// NewAPI returns the API for the given host.
//
// In the future some kind of validation might be put here that the
// host is compatible with the given implementation.
func NewAPI(host string, settings ...ant.Setting) *API {
	api := API{
		host:   host,
		client: http.DefaultClient,
	}
	api.SetCredentials(nil)
	ant.MustConfigure(&api, settings...)
	return &api
}

// API provides http request builders for the tidio service
// The requests returned should be valid and complete.
type API struct {
	host   string
	client *http.Client
	auth   ant.Setting // applied

	// last api
	Request *http.Request
}

func (me *API) CreateTimesheet(path string, body io.Reader) *API {
	r := me.newRequest("POST", path, body)
	me.Auth(r)
	return me
}

func (me *API) ReadTimesheet(path string) *API {
	r := me.newRequest("GET", path, nil)
	me.Auth(r)
	return me
}

func (me *API) SetCredentials(c *Credentials) {
	me.auth = NewBasicAuth(c)
}

// Auth applies credentials to the request and sets them as last
// values on the api.
func (me *API) Auth(r *http.Request) {
	if r == nil {
		return
	}
	err := ant.Configure(r, me.auth)
	me.warn(err)
}

// newRequest
func (me *API) newRequest(method, path string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, me.host+path, body)
	me.warn(err)
	me.Request = r
	return r
}

// MustSend
func (me *API) MustSend() *http.Response {
	r, err := me.Send()
	if err != nil {
		panic(err)
	}
	return r
}

func (me *API) Send() (*http.Response, error) {
	r := me.Request
	if r == nil {
		return nil, fmt.Errorf("missing API request")
	}
	resp, err := me.client.Do(r)
	if err != nil {
		Log(me).Error(r.Method, r.URL, err)
		return resp, err
	}
	Log(me).Info(r.Method, r.URL, resp.StatusCode)
	return resp, nil
}

// warn logs non nil errors
func (me *API) warn(err error) {
	if err == nil {
		return
	}
	Log(me).Info(err)
}

// ----------------------------------------

func NewConfig() *Config {
	c := &Config{
		activeLoggers: make(map[interface{}]*LogPrinter),
	}
	c.SetOutput(ioutil.Discard)
	return c
}

// Config holds reference to loggers for various objects.
type Config struct {
	out io.Writer

	debug bool

	activeLoggers map[interface{}]*LogPrinter
}

// SetOutput
func (me *Config) SetOutput(w io.Writer) { me.out = w }

// Debug
func (me *Config) Debug() bool { return me.debug }
func (me *Config) SetDebug(v bool) {
	me.debug = v
}

// Unregister removes the previously registered item if any.
func (me *Config) Unregister(v interface{}) {
	delete(me.activeLoggers, v)
}

func (me *Config) Log(v interface{}) *LogPrinter {
	l, found := me.activeLoggers[v]
	if !found {
		return nolog
	}
	return l
}

// Register creates a LogPrinter for the given objects.
// The registered LogPrinter is later retrieved with Config.Log(v)
func (me *Config) Register(v ...interface{}) *LogPrinter {
	if len(v) == 0 {
		panic("missing values in MLog")
	}

	first := v[0]
	l := NewLogPrinter(me.out)
	me.activeLoggers[first] = l
	for _, other := range v[1:] {
		me.activeLoggers[other] = l
	}
	return l
}

// ----------------------------------------

// Wrappers for default config Conf

func Register(v ...interface{}) *LogPrinter {
	return Conf.Register(v...)
}
func Log(v interface{}) *LogPrinter { return Conf.Log(v) }
func Unregister(v interface{})      { Conf.Unregister(v) }

var Conf = NewConfig()

// ----------------------------------------

var nolog = &LogPrinter{
	lgr: log.New(ioutil.Discard, "", 0),
}

// ----------------------------------------

func NewHelpView() *Page {
	nav := Nav()
	content := Article(
		NewAPISection(),
		Section(
			H2("Timesheet file format"),
			P("Timesheets are plain text and are specific to year and month"),
			Pre(Class("timesheet"), timesheet201506),
		),

		NewChangelog(),
	)
	body := Body(
		Header(
			H1("Tidio - API documentation"),
		),
		nav,
		content,
		footer(),
	)
	toc.MakeTOC(nav, body, "h1", "h2", "h3")
	return NewPage(
		Html(
			Head(
				Title("tidio - help"),
				//apidoc.DefaultStyle(),
				Style(theme()),
			),
			body,
		),
	)
}

func NewChangelog() *Element {
	return Article(
		H1("Changelog"),

		Pre(changelog),
	)
}

// Version returns the latest version according to the embeded
// changelog.
func Version() string {
	from := strings.Index(changelog, "[")
	to := strings.Index(changelog, "]")
	return changelog[from+1 : to]
}

//go:embed changelog.md
var changelog string

func NewAPISection() *Element {
	// Cache api section
	cred := NewCredentials("john", "secret")
	sys := NewSystem(cred)
	doc := apidoc.NewDoc(NewRouter(sys))
	api := NewAPI("https://tidio.preferit.se")
	ant.MustConfigure(api, cred)

	return Section(
		H2("Timesheets"),
		P(
			``,
		),
		H3("Create or update"),
		doc.Use(api.CreateTimesheet(
			"/api/timesheets/john/201506.timesheet",
			strings.NewReader(timesheet201506),
		).Request),
		doc.JsonResponse(),

		H3("Read specific timesheet"),
		doc.Use(api.ReadTimesheet(
			"/api/timesheets/john/201506.timesheet",
		).Request),
		doc.Response(),
	)
}

const timesheet201506 = `2015 June
---------
23  1 Mon 8
    2 Tue 8
    3 Wed 8 (3 meeting)
    4 Thu 8
    5 Fri 6 Ended work 2 hours early, felt sick.
    6 Sat
    7 Sun
24  8 Mon 8
    9 Tue 8
   10 Wed 8
   11 Thu 8 (7 conference) (1 travel)
   12 Fri 8
   13 Sat
   14 Sun
25 15 Mon 8
   16 Tue 8
   17 Wed 8:30
   18 Thu 8
   19 Fri 8
   20 Sat
   21 Sun
26 22 Mon 8
   23 Tue 8
   24 Wed 8
   25 Thu 8
   26 Fri 8
   27 Sat
   28 Sun
27 29 Mon 8
   30 Tue 8`

func theme() *CSS {
	css := NewCSS()
	css.Style("html, body",
		"margin: 0 0",
		"padding: 0 0",
		"background-color: #ffffff",
	)
	css.Style("h1:first-child",
		"margin-top: 0",
	)
	css.Style("a:link",
		"color: rgb(55, 94, 171)", // golang blue
		"text-decoration: none",
	)
	css.Style("a:link:hover",
		"text-decoration: underline",
	)
	css.Style("header",
		"padding-top: 1em",
		"padding-left: 1.62em",
	)
	css.Style("nav",
		"padding-left: 1.62em",
		"font-family: Arial, Helvetica, sans-serif",
	)
	css.Style("article",
		"background-color: white",
		"padding: 1em 1em 2em 1.62em",
		"min-height: 300",
	)
	css.Style("section",
		"margin-bottom: 1.62em",
	)
	css.Style("pre",
		"margin-left: 1.62em",
	)
	css.Style("footer",
		"border-top: 1px solid #727272",
		"padding: 0.6em 0.6em",
		"background-color: #e2e2e2",
		"min-height: 500px",
	)
	css.Style(".timesheet",
		"border: 1px #e2e2e2 dotted",
		"padding: 1em 1em",
		"background-color: #ffffe6",
	)
	css.Style("p",
		"font-family: Arial, Helvetica, sans-serif",
		"line-height: 1.3em",
	)
	css.Style(".request",
		"padding: 1em 1.618em",
		"border-radius: 1em",
		"border: 1px dashed #929292",
	)
	css.Style(".response",
		"padding: 1em 1.618em",
		"background-color: #f2f2f2",
		"border-radius: 1em",
	)
	css.Style("nav ul",
		"list-style-type: none",
		"padding-left: 0",
		"line-height: 1.3em",
	)
	css.Style("nav ul .h2",
		"margin-left: 1.62em",
	)
	css.Style("nav ul .h3",
		"margin-left: 3.22em",
	)

	return css
}

// When the System started so we know the uptime
var serviceStarted = time.Now()

func footer() *Element {
	return Footer(
		"Uptime: ",
		time.Since(serviceStarted).Round(time.Second).String(),
	)
}

// ----------------------------------------

func NewLogPrinter(w io.Writer) *LogPrinter {
	return &LogPrinter{
		lgr: log.New(w, "", log.Lshortfile),
	}
}

type LogPrinter struct {
	buf    bytes.Buffer // if buffered
	lgr    *log.Logger
	writes int
	failed bool
}

// Buf makes the log printer buffered. Use Flush to get the contents.
func (me *LogPrinter) Buf() *LogPrinter {
	me.SetOutput(&me.buf)
	return me
}

// FlushString
func (me *LogPrinter) FlushString() string {
	return string(me.Flush())
}

// Flush returns the buffered bytes if any and resets the buffer.
func (me *LogPrinter) Flush() []byte {
	defer me.buf.Reset()
	return me.buf.Bytes()
}

// Failed
func (me *LogPrinter) Failed() bool { return me.failed }

// Info
func (me *LogPrinter) Info(v ...interface{}) {
	me.lgr.Output(2, fmt.Sprintln(v...))
	me.writes++
}

// Error
func (me *LogPrinter) Error(v ...interface{}) {
	me.lgr.Output(2, fmt.Sprintln(v...))
	me.writes++
	me.failed = true
}

func (me *LogPrinter) Log(v ...interface{}) {
	me.lgr.Output(2, fmt.Sprintln(v...))
	me.writes++
}

func (me *LogPrinter) SetOutput(w io.Writer) { me.lgr.SetOutput(w) }
