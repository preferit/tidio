package tidio

import (
	"strings"
	"time"

	"github.com/gregoryv/ant"
	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/apidoc"
)

func NewHelpView() *Page {
	john := &InitialAccount{account: "john", secret: "secret"}
	srv := NewService(john)
	doc := apidoc.NewDoc(srv.Router())
	api := NewAPI("https://tidio.preferit.se")
	asJohn := NewCredentials(john.account, john.secret)
	ant.MustConfigure(api, asJohn)

	content := Div(
		Section(
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
		),

		Section(
			H2("Timesheet file format"),
			P("Timesheets are plain text and are specific to year and month"),
			Pre(Class("timesheet"), timesheet201506),
		),
	)

	return NewPage(
		"help.html",
		Html(
			Head(
				apidoc.DefaultStyle,
				Style(theme()),
			),
			Body(content, footer()),
		),
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
		"background-color: #e2e2e2",
	)
	css.Style("div",
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

// When the service started so we know the uptime
var serviceStarted = time.Now()

func footer() *Element {
	return Footer(
		"Uptime: ",
		time.Since(serviceStarted).Round(time.Second).String(),
	)
}
