// Command tidio is a standalone http server for the tidio.Service
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/fox"
	"github.com/preferit/tidio"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	c := NewCommand()
	c.Logger = fox.NewSyncLog(os.Stderr).FilterEmpty()
	if err := c.run(os.Args...); err != nil {
		c.Log(err)
		os.Exit(1)
	}
}

// NewCommand returns a command without logging and default options
func NewCommand() *Command {
	return &Command{
		ListenAndServe: http.ListenAndServe,
		Logger:         fox.NewSyncLog(ioutil.Discard),
	}
}

type Command struct {
	ListenAndServe func(string, http.Handler) error
	fox.Logger
}

func (c *Command) run(args ...string) error {
	cli := cmdline.New(args...)

	// parse arguments
	bind, opt := cli.Option("-bind").StringOpt(":13001")
	opt.Doc("[host]:port to bind to")

	stateFilename, opt := cli.Option("-state").StringOpt("system.state")
	opt.Doc("file to keep state in")
	if stateFilename == "" {
		return fmt.Errorf("-state cannot be empty")
	}

	if err := cli.Error(); err != nil {
		var buf bytes.Buffer
		cli.WriteUsageTo(&buf)
		c.Log(buf.String())
		return err
	}

	service := tidio.NewService()
	service.SetLogger(c.Logger)

	if err := c.initStateRestoration(service, stateFilename); err != nil {
		return err
	}

	c.Log("listen on ", bind)
	return c.ListenAndServe(bind, service)
}

// initStateRestoration
func (me *Command) initStateRestoration(service *tidio.Service, filename string) error {
	dest := tidio.NewFileStorage(filename)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := service.PersistState(dest); err != nil {
			return err
		}
	} else {
		if err := service.RestoreState(filename); err != nil {
			return err
		}
	}
	service.AutoPersist(dest)
	return nil
}
