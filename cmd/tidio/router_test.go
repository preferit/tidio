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
	apikeys := tidio.APIKeys{
		"KEY": "john",
	}
	headers := http.Header{}
	store, cleanup := newTempStore(t)
	defer cleanup()
	service := tidio.NewService(store, apikeys)
	router := NewRouter(service)

	exp := assert().ResponseFrom(router)
	exp.StatusCode(200, "GET", "/api", nil)
	exp.Contains("revision", "GET", "/api")
	exp.Contains("version", "GET", "/api")
	exp.Contains("resources", "GET", "/api")

	headers = http.Header{}
	headers.Set("Authorization", "NO SUCH KEY")
	exp.StatusCode(401, "POST", "/api/timesheets/not_there/202001.timesheet",
		strings.NewReader("body"), headers)

	// authenticated
	headers = http.Header{}
	headers.Set("Authorization", "KEY")
	exp.StatusCode(404, "GET", "/api/timesheets/", headers)
	exp.StatusCode(204, "POST", "/api/timesheets/john/202001.timesheet",
		strings.NewReader("some content"), headers)
	exp.StatusCode(403, "POST", "/api/timesheets/eva/199601.timesheet",
		strings.NewReader("body"), headers)
	exp.StatusCode(400, "POST", "/api/timesheets/john/202001xtimesheet",
		strings.NewReader("body"), headers)

	// read timesheet
	exp.StatusCode(400, "GET", "/api/timesheets/john/999900.timesheet", headers)
	content := "TEST content"
	exp.StatusCode(204, "POST", "/api/timesheets/john/197604.timesheet",
		strings.NewReader(content), headers)
	exp.StatusCode(200, "GET", "/api/timesheets/john/197604.timesheet", headers)
	exp.Contains(content, "GET", "/api/timesheets/john/197604.timesheet", headers)

}

func Test_convert_error(t *testing.T) {
	ok := func(err error, exp int) {
		t.Helper()
		got := statusOf(err)
		if got != exp {
			t.Error("got:", got, "exp:", exp)
		}
	}
	ok(nil, http.StatusOK)
	ok(tidio.ErrForbidden, http.StatusForbidden)
}
