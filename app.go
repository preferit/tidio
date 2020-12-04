package tidio

import (
	"flag"
	"net/http"
	"os"

	"github.com/gregoryv/ant"
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
		fs       = flag.NewFlagSet(cmd.Args()[0], flag.ContinueOnError)
		bind     = fs.String("bind", ":13001", "[host]:port to bind to")
		filename = fs.String("state", "system.state", "")
	)
	fs.SetOutput(cmd.Stderr())
	err := fs.Parse(cmd.Args()[1:])
	if err != nil {
		if err != flag.ErrHelp {
			return err
		}
		return nil
	}

	srv := NewService()
	ant.MustConfigure(srv, fox.Logging{me})

	srv.dest = NewFileStorage(*filename) // todo replace with ant Setting

	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		if err := srv.PersistState(); err != nil {
			return err
		}
	} else {
		if err := srv.RestoreState(*filename); err != nil {
			return err
		}
	}
	me.Log("listening on:", *bind)
	return me.ListenAndServe(*bind, srv.Router())
}
