package tidio

import (
	"net/http"
	"os"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/fox"
	"github.com/gregoryv/wolf"
)

// NewApp returns a App with applied settings
func NewApp(settings ...ant.Setting) *App {
	app := App{
		ListenAndServe: http.ListenAndServe,
	}
	ant.MustConfigure(&app, settings...)
	return &app
}

type App struct {
	OptionalLogger
	wolf.Command
	ListenAndServe func(string, http.Handler) error
}

func (me *App) Run(cmd wolf.Command) int {
	if err := me.run(cmd); err != nil {
		me.Log(err)
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

// ----------------------------------------

type serveHTTP struct {
	cmdline.Item
	bind     string
	filename string
}

// ExtraOptions
func (me *serveHTTP) ExtraOptions(p *cmdline.Parser) {
	me.bind = p.Option("-bind").String(":13001")
	me.filename = p.Option("-state").String("system.state")
}

func (me *serveHTTP) Run(app *App) error {
	srv := NewService()
	ant.MustConfigure(srv, fox.Logging{app})

	// configure persistence
	srv.dest = NewFileStorage(me.filename)

	if _, err := os.Stat(me.filename); os.IsNotExist(err) {
		if err := srv.PersistState(); err != nil {
			return err
		}
	} else {
		if err := srv.RestoreState(); err != nil {
			return err
		}
	}
	app.Log("listening on:", me.bind)
	return app.ListenAndServe(me.bind, srv.Router())
}

// ----------------------------------------

type runnable interface {
	Run(app *App) error
}
