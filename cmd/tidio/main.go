// Command tidio is a standalone http server for the tidio.Service
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/stamp"
	"github.com/preferit/tidio"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	c := &cli{
		ListenAndServe: http.ListenAndServe,
		Logger:         fox.NewSyncLog(os.Stderr).FilterEmpty(),
	}
	stamp.InitFlags()
	flag.StringVar(&c.bind, "bind", ":13001", "[host]:port to bind to")
	flag.StringVar(&c.stateFilename, "state", "system.state", "file to keep state in")
	flag.Parse()
	stamp.AsFlagged()

	if err := run(c); err != nil {
		c.Log(err)
		os.Exit(1)
	}
}

type cli struct {
	ListenAndServe func(string, http.Handler) error
	bind           string
	stateFilename  string
	fox.Logger
}

func run(c *cli) error {
	if c.bind == "" {
		return fmt.Errorf("empty bind")
	}
	service := tidio.NewService()
	service.SetLogger(c.Logger)

	if c.stateFilename != "" {
		c.initStateRestoration(service)
	}
	c.Log("listen on ", c.bind)
	return c.ListenAndServe(c.bind, service)
}

// initStateRestoration
func (me *cli) initStateRestoration(service *tidio.Service) {
	dest := tidio.NewFileStorage(me.stateFilename)
	if _, err := os.Stat(me.stateFilename); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		me.Log("creating ", path.Join(wd, me.stateFilename))
		if err := service.PersistState(dest); err != nil {
			me.Log(err)
		}

	} else {
		if err := service.RestoreState(me.stateFilename); err != nil {
			me.Log(err)
		}
	}
	service.AutoPersist(dest)
}
