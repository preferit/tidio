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
	os.MkdirAll(c.out, 0755)
	indexPage().SaveTo(c.out)
	return nil
}

func indexPage() *Page {
	article := Article(
		H1("Tidio"),
		A(Href("/api"), "/api"),
		H2("Timesheets"),
		H3("Create or update"),
		Pre(`HTTP/1.1 POST {host}/api/timesheets/{account}/{yyyymm}.timesheet
Authorization: {key}

-- body contains timesheet --`),

		H3("Read specific timesheet"),
		Pre(`HTTP/1.1 GET {host}/api/timesheets/{account}/{yyyymm}.timesheet
Authorization: {key}`),
		"Responds with timesheet",

		H3("List timesheets of a specific user"),
		Pre(`HTTP/1.1 GET {host}/api/timesheets/{account}/
Authorization: {key}`),
		`Responds with json {"timesheets": []}`,
	)

	return NewPage(
		"index.html",
		Html(
			Head(
				Style(theme()),
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

func theme() *CSS {
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
	css.Style("pre",
		"margin-left: 1em",
	)
	css.Style("footer",
		"border-top: 1px solid #727272",
		"padding: 0.6em 0.6em",
	)
	return css
}
