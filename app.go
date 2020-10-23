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

// NewApp returns a App without logging and default options
func NewApp(cmd wolf.Command, settings ...ant.Setting) *App {
	app := App{
		Command:        cmd,
		ListenAndServe: http.ListenAndServe,
		Logger: Logger{
			fox.NewSyncLog(cmd.Stderr()).FilterEmpty(),
		},
	}
	ant.MustConfigure(&app, settings...)
	return &app
}

type App struct {
	Logger
	wolf.Command
	ListenAndServe func(string, http.Handler) error
}

func (me *App) Run() int {
	if err := me.run(); err != nil {
		me.Log(err)
		return me.Stop(1)
	}
	return me.Stop(0)
}

func (me *App) run() error {
	var (
		fs            = flag.NewFlagSet(me.Args()[0], flag.ContinueOnError)
		bind          = fs.String("bind", ":13001", "[host]:port to bind to")
		stateFilename = fs.String("state", "system.state", "")
	)
	fs.SetOutput(me.Stderr())
	err := fs.Parse(me.Args()[1:])
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
	me.Log("add account john")
	me.Log(srv.AddAccount("john", "secret"))
	me.Log(srv.InitResources())
	me.Log("listen on ", *bind)
	return me.ListenAndServe(*bind, srv)
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
	service.AutoPersist(dest, 3*time.Second)
	return nil
}
