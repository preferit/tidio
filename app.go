package tidio

import (
	"flag"
	"net/http"
	"os"
	"time"

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
		fs            = flag.NewFlagSet(cmd.Args()[0], flag.ContinueOnError)
		bind          = fs.String("bind", ":13001", "[host]:port to bind to")
		stateFilename = fs.String("state", "system.state", "")
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

	if err := me.initStateRestoration(srv, *stateFilename); err != nil {
		return err
	}
	return me.ListenAndServe(*bind, srv.Router())
}

// initStateRestoration
func (me *App) initStateRestoration(service *Service, filename string) error {
	dest := NewFileStorage(filename)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := service.PersistState(dest); err != nil {
			return err
		}
	} else {
		if err := service.RestoreState(filename); err != nil {
			return err
		}
	}
	// todo replace this with something else
	service.AutoPersist(dest, 3*time.Second)
	return nil
}
