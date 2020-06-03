package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/preferit/tidio"
)

func Test_router(t *testing.T) {
	assert := asserter.New(t)
	apikeys := map[string]string{
		"KEY": "john",
	}
	headers := http.Header{}
	store, cleanup := newTempStore(t)
	defer cleanup()

	service := tidio.NewService(store, apikeys)
	router := NewRouter(apikeys, store, service)
	exp := assert().ResponseFrom(router)
	exp.StatusCode(200, "GET", "/api", nil)
	exp.Contains("revision", "GET", "/api")
	exp.Contains("version", "GET", "/api")
	exp.Contains("resources", "GET", "/api")

	exp.StatusCode(401, "GET", "/api/timesheets/")
	headers = http.Header{}
	headers.Set("Authorization", "NO SUCH KEY")
	exp.StatusCode(401, "GET", "/api/timesheets/", headers)
	exp.StatusCode(401, "POST", "/api/timesheets/not_there/202001.timesheet",
		strings.NewReader("body"), headers)

	// authenticated
	headers = http.Header{}
	headers.Set("Authorization", "KEY")
	exp.StatusCode(200, "GET", "/api/timesheets/", headers)
	exp.StatusCode(204, "POST", "/api/timesheets/john/202001.timesheet",
		strings.NewReader("some content"), headers)
	exp.StatusCode(403, "POST", "/api/timesheets/eva/199601.timesheet",
		strings.NewReader("body"), headers)
	exp.StatusCode(400, "POST", "/api/timesheets/john/202001xtimesheet",
		strings.NewReader("body"), headers)
}
