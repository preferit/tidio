package tidio

import (
	"strings"
	"time"

	"github.com/gregoryv/ant"
	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/apidoc"
	"github.com/gregoryv/web/toc"
)

func NewHelpView() *Page {
	nav := Nav()
	content := Article(
		apiSection,
		Section(
			H2("Timesheet file format"),
			P("Timesheets are plain text and are specific to year and month"),
			Pre(Class("timesheet"), timesheet201506),
		),
	)
	toc.MakeTOC(nav, content, "h2", "h3")
	return NewPage(
		Html(
			Head(
				Title("tidio - help"),
				apidoc.DefaultStyle(),
				Style(theme()),
			),
			Body(
				Header(
					H1("Tidio - API documentation"),
				),
				nav,
				content,
				footer(),
			),
		),
	)
}

var apiSection *Element

func init() {
	// Cache api section
	cred := NewCredentials("john", "secret")
	srv := NewService(cred)
	doc := apidoc.NewDoc(srv.Router())
	api := NewAPI("https://tidio.preferit.se")
	ant.MustConfigure(api, cred)

	apiSection = Section(
		H2("Timesheets"),
		P(
			``,
		),
		H3("Create or update"),
		doc.Use(api.CreateTimesheet(
			"/api/timesheets/john/201506.timesheet",
			strings.NewReader(timesheet201506),
		).Request),
		doc.JsonResponse(),

		H3("Read specific timesheet"),
		doc.Use(api.ReadTimesheet(
			"/api/timesheets/john/201506.timesheet",
		).Request),
		doc.Response(),
	)
}

const timesheet201506 = `2015 June
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
   30 Tue 8`

func theme() *CSS {
	css := NewCSS()
	css.Style("html, body",
		"margin: 0 0",
		"padding: 0 0",
		"background-color: #ffffff",
	)
	css.Style("a:link",
		"color: rgb(55, 94, 171)", // golang blue
		"text-decoration: none",
	)
	css.Style("a:link:hover",
		"text-decoration: underline",
	)
	css.Style("header",
		"padding-top: 1em",
		"padding-left: 1.62em",
	)
	css.Style("nav",
		"padding-left: 1.62em",
		"font-family: Arial, Helvetica, sans-serif",
	)
	css.Style("article",
		"background-color: white",
		"padding: 1em 1em 2em 1.62em",
		"min-height: 300",
	)
	css.Style("section",
		"margin-bottom: 1.62em",
	)
	css.Style("pre",
		"margin-left: 1.62em",
	)
	css.Style("footer",
		"border-top: 1px solid #727272",
		"padding: 0.6em 0.6em",
		"background-color: #e2e2e2",
	)
	css.Style(".timesheet",
		"border: 1px #e2e2e2 dotted",
		"padding: 1em 1em",
		"background-color: #ffffe6",
	)
	css.Style("p",
		"font-family: Arial, Helvetica, sans-serif",
	)
	return css
}

// When the service started so we know the uptime
var serviceStarted = time.Now()

func footer() *Element {
	return Footer(
		"Uptime: ",
		time.Since(serviceStarted).Round(time.Second).String(),
	)
}
