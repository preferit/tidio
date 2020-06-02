package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	. "github.com/gregoryv/web"
)

func main() {
	c := &cli{}
	flag.StringVar(&c.out, "o", "/tmp", "root of site")
	flag.Parse()

	if err := c.run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cli struct {
	out string
}

func (c *cli) run() error {
	indexPage().SaveTo(c.out)
	return nil
}

func indexPage() *Page {
	article := Article(
		H1("Tidio"),
		A(Href("/api"), "/api"),
		Br(),
		time.Now().String(),
	)
	return NewPage("index.html", Html(Body(article)))
}
