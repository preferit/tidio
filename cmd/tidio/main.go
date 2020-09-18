// Command tidio is a standalone http server for the tidio.Service
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

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
	flag.StringVar(&c.stateFilename, "state", "", "file to keep state in")
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

	if err := service.RestoreState(c.stateFilename); err != nil {
		return err
	}

	c.Log("listen on ", c.bind)
	return c.ListenAndServe(c.bind, service)
}
