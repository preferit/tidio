package tidio

import (
	"net/http"
	"os"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/wolf"
)

// NewApp returns a App with applied settings
func NewApp(settings ...ant.Setting) *App {
	app := &App{
		ListenAndServe: http.ListenAndServe,
	}
	ant.MustConfigure(app, settings...)
	RLog(app)
	return app
}

type App struct {
	wolf.Command
	ListenAndServe func(string, http.Handler) error
}

func (me *App) Run(cmd wolf.Command) int {
	if err := me.run(cmd); err != nil {
		Log(me).Info("Run", err)
		return cmd.Stop(1)
	}
	return cmd.Stop(0)
}

func (me *App) run(cmd wolf.Command) error {
	var (
		cli  = cmdline.NewParser(cmd.Args()...)
		help = cli.Flag("-h, --help")

		a      = cli.Group("Actions", "ACTION")
		_      = a.New("serveHTTP", &serveHTTP{})
		action = a.Selected()
	)
	switch {
	case !cli.Ok():
		return cli.Error()

	case help:
		cli.WriteUsageTo(cmd.Stdout())
		return nil
	}
	return action.(runnable).Run(me)

}

type runnable interface {
	Run(app *App) error
}

// ----------------------------------------

type serveHTTP struct {
	cmdline.Item
	bind     string
	filename string
}

func (me *serveHTTP) ExtraOptions(p *cmdline.Parser) {
	me.bind = p.Option("-bind").String(":13001")
	me.filename = p.Option("-state").String("system.state")
}

func (me *serveHTTP) Run(app *App) error {
	srv := NewService()

	// configure persistence
	srv.dest = NewFileStorage(me.filename)

	_, err := os.Stat(me.filename)
	switch {
	case os.IsNotExist(err):
		srv.PersistState()
	default:
		srv.RestoreState()
	}
	if err := srv.Error(); err != nil {
		return err
	}
	Log(app).Info("listening on", me.bind)
	return app.ListenAndServe(me.bind, srv.Router())
}

// ----------------------------------------
