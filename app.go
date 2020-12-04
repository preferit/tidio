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
		cli      = cmdline.NewParser(cmd.Args()...)
		bind     = cli.Option("-bind").String(":13001")
		filename = cli.Option("-state").String("system.state")
		help     = cli.Flag("-h, --help")
	)
	switch {
	case !cli.Ok():
		return cli.Error()

	case help:
		cli.WriteUsageTo(cmd.Stdout())
		return nil
	}

	srv := NewService()
	ant.MustConfigure(srv, fox.Logging{me})

	// configure persistence
	srv.dest = NewFileStorage(filename)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := srv.PersistState(); err != nil {
			return err
		}
	} else {
		if err := srv.RestoreState(filename); err != nil {
			return err
		}
	}
	me.Log("listening on:", bind)
	return me.ListenAndServe(bind, srv.Router())
}
