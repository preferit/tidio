package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gregoryv/stamp"
	"github.com/preferit/tidio"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	c := &cli{
		starter: http.ListenAndServe,
	}
	stamp.InitFlags()
	flag.StringVar(&c.storeDir, "store-dir", "./store", "file storage repository")
	flag.StringVar(&c.keysfile, "keys", "apikeys.json", "map of apikeys")
	flag.StringVar(&c.bind, "bind", ":13001", "[host]:port to bind to")
	flag.Parse()
	stamp.AsFlagged()

	if err := c.run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cli struct {
	storeDir string
	bind     string
	keysfile string
	starter  func(string, http.Handler) error
}

func (c *cli) run() error {
	if c.bind == "" {
		return fmt.Errorf("empty bind")
	}
	accounts := tidio.AccountsMap{}.New()
	if err := accounts.ReadState(func() (io.ReadCloser, error) {
		return os.Open(c.keysfile)
	}); err != nil {
		return err
	}

	os.MkdirAll(c.storeDir, 0755)
	service := tidio.Service{}.New()
	service.Accounts = accounts
	service.LoadState(path.Join(c.storeDir, "data.gob"))
	service.SaveState() // update
	router := NewRouter(service)
	return c.starter(c.bind, router)
}
