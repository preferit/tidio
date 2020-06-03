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
	)
	css := NewCSS()
	css.Style("html, body",
		"margin: 0 0",
		"padding: 0 0",
		"background-color: #e2e2e2",
	)
	css.Style("article",
		"background-color: white",
		"padding: 1em 1em 2em 1em",
		"min-height: 300",
	)
	css.Style("footer",
		"border-top: 1px solid #727272",
		"padding: 0.6em 0.6em",
	)

	return NewPage(
		"index.html",
		Html(
			Head(
				Style(css),
			),
			Body(article,

				Footer(
					"Generated: ",
					time.Now().Round(time.Second).String(),
				),
			),
		),
	)
}
