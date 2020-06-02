package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gregoryv/stamp"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	c := &cli{
		starter: http.ListenAndServe,
	}
	stamp.InitFlags()
	flag.StringVar(&c.bind, "bind", ":13001", "[host]:port to bind to")
	flag.Parse()
	stamp.AsFlagged()

	if err := c.run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cli struct {
	bind    string
	starter func(string, http.Handler) error
}

func (c *cli) run() error {
	if c.bind == "" {
		return fmt.Errorf("empty bind")
	}
	router := NewRouter()
	return c.starter(c.bind, router)
}
