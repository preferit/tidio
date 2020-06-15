package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gregoryv/stamp"
	"github.com/preferit/tidio"
)

//go:generate stamp -clfile ../../changelog.md -go build_stamp.go

func main() {
	c := &cli{
		ListenAndServe: http.ListenAndServe,
	}
	stamp.InitFlags()
	flag.StringVar(&c.bind, "bind", ":13001", "[host]:port to bind to")
	flag.Parse()
	stamp.AsFlagged()

	if err := run(c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cli struct {
	ListenAndServe func(string, http.Handler) error
	bind           string
}

func run(c *cli) error {
	if c.bind == "" {
		return fmt.Errorf("empty bind")
	}
	service := tidio.NewService()
	return c.ListenAndServe(c.bind, service)
}
