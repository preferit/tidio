package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	c := &cli{
		starter: http.ListenAndServe,
	}
	flag.StringVar(&c.bind, "bind", ":13001", "[host]:port to bind to")
	flag.Parse()
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "tidio API")
	})
	return c.starter(c.bind, nil)
}
