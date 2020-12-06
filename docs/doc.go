package docs

import (
	"github.com/gregoryv/go-timesheet"
	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/apidoc"
	"github.com/preferit/tidio"
)

func NewIndex() *Page {
	var (
		asJohn = tidio.NewCredentials("john", "secret")
		api    = tidio.NewAPI("", asJohn)

		initAcc = asJohn
		srv     = tidio.NewService(initAcc)

		router = srv.Router()
		doc    = apidoc.NewDoc(router)
	)

	page := NewPage(Html(Body(
		doc.Use(api.CreateTimesheet(
			"/api/timesheets/john/202001.timesheet",
			timesheet.Render(2020, 1, 8),
		).Request),
		doc.JsonResponse(),
	)))
	return page
}
