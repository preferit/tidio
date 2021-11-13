// Tidio tidio is a standalone http server for the tidio.Service
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/rs"
	"github.com/preferit/tidio"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

// By default called from a os shell
var cmd = cmdline.NewShellOS()

func main() {
	conf := tidio.Conf
	conf.SetOutput(os.Stderr)

	app := NewApp()

	var (
		cli = cmdline.NewBasicParser()

		a      = cli.Group("Actions", "ACTION")
		_      = a.New("serveHTTP", &serveHTTP{})
		_      = a.New("mkAccount", &mkAccount{})
		action = a.Selected()
	)
	cli.Parse()

	action.(runnable).Run(app)
}

type runnable interface {
	Run(app *App) error
}

// NewApp returns a App with applied settings
func NewApp() *App {
	app := &App{
		ListenAndServe: http.ListenAndServe,
	}
	tidio.Register(app)
	return app
}

type App struct {
	ListenAndServe func(string, http.Handler) error
}

// ----------------------------------------

type serveHTTP struct {
	cmdline.Item
	bind     string
	filename string
}

func (me *serveHTTP) ExtraOptions(p *cmdline.Parser) {
	me.bind = p.Option("-bind").String(":13001")
	me.filename = p.Option("-state").String("/var/local/tidio/system.state")
}

func (me *serveHTTP) Run(app *App) error {
	srv := tidio.NewService()

	// configure persistence
	srv.UseFileStorage(me.filename)

	tidio.Log(app).Info("listening on", me.bind)
	return app.ListenAndServe(me.bind, srv.Router())
}

// ----------------------------------------

type mkAccount struct {
	cmdline.Item

	filename string
	name     string
	secret   string

	writeBack bool
}

func (me *mkAccount) ExtraOptions(p *cmdline.Parser) {
	me.filename = p.Option("-state").String("/var/local/tidio/system.state")
	me.name = p.Option("-n, --name").String("")
	me.secret = p.Option("-s, --secret").String("")
	me.writeBack = p.Flag("-w, --write")
}

func (me *mkAccount) Run(app *App) error {
	r, err := os.Open(me.filename)
	if err != nil {
		return fmt.Errorf("open state file: %w", err)
	}
	defer r.Close()
	sys := rs.NewSystem()
	err = sys.Import("/", r)
	if err != nil {
		return err
	}

	asRoot := rs.Root.Use(sys)

	sh := tidio.NewShell(asRoot)
	err = tidio.CreateAccount(sh, me.name, me.secret)
	if err != nil {
		return err
	}
	out := cmd.Stdout()
	if me.writeBack {
		out, err := os.Create(me.filename)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return sys.Export(out)
}
