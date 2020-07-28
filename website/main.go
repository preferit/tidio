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
		H1("tidio"),
		A(Href("/api"), "/api"),

		Section(
			H2("Timesheets"),
			P(
				``,
			),
			H3("Create or update"),
			Pre(`HTTP/1.1 POST {host}/api/timesheets/{account}/{yyyymm}.timesheet
Authorization: {auth}

-- body contains timesheet --`),

			H3("Read specific timesheet"),
			Pre(`HTTP/1.1 GET {host}/api/timesheets/{account}/{yyyymm}.timesheet
Authorization: {auth}`),
			"Responds with timesheet",

			H3("List timesheets of a specific user"),
			Pre(`HTTP/1.1 GET {host}/api/timesheets/{account}/
Authorization: {auth}`),
			`Responds with json {"timesheets": []}`,
		),

		Section(
			H2("Timesheet file format"),
			P("Timesheets are plain text and are specific to year and month"),
			Pre(Class("timesheet"),
				`2015 June
---------
23  1 Mon 8
    2 Tue 8
    3 Wed 8 (3 meeting)
    4 Thu 8
    5 Fri 6 Ended work 2 hours early, felt sick.
    6 Sat
    7 Sun
24  8 Mon 8
    9 Tue 8
   10 Wed 8
   11 Thu 8 (7 conference) (1 travel)
   12 Fri 8
   13 Sat
   14 Sun
25 15 Mon 8
   16 Tue 8
   17 Wed 8:30
   18 Thu 8
   19 Fri 8
   20 Sat
   21 Sun
26 22 Mon 8
   23 Tue 8
   24 Wed 8
   25 Thu 8
   26 Fri 8
   27 Sat
   28 Sun
27 29 Mon 8
   30 Tue 8`,
			),
		),
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
	css.Style("section",
		"margin-bottom: 5em",
	)
	css.Style("pre",
		"margin-left: 1em",
	)
	css.Style("footer",
		"border-top: 1px solid #727272",
		"padding: 0.6em 0.6em",
	)
	css.Style(".timesheet",
		"border: 1px #e2e2e2 dotted",
		"padding: 1em 1em",
		"background-color: #ffffe6",
	)
	return css
}
