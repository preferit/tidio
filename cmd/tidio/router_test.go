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
		"KEY": tidio.NewAccount("john", "admin"),
	}
	headers := http.Header{}
	service := tidio.NewService(apikeys)
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
	exp.StatusCode(204, "POST", "/api/timesheets/eva/199601.timesheet",
		strings.NewReader("body"), headers)
	exp.StatusCode(400, "POST", "/api/timesheets/john/202001xtimesheet",
		strings.NewReader("body"), headers)

	// read timesheet
	exp.StatusCode(400, "GET", "/api/timesheets/john/999900.timesheet", headers)
	exp.StatusCode(204, "POST", "/api/timesheets/john/197604.timesheet",
		strings.NewReader(timesheet197604), headers)
	exp.StatusCode(200, "GET", "/api/timesheets/john/197604.timesheet", headers)
	exp.BodyIs(timesheet197604, "GET", "/api/timesheets/john/197604.timesheet", headers)
	exp.Contains("197604.timesheet", "GET", "/api/timesheets/john/", headers)
	exp.StatusCode(200, "GET", "/api/timesheets/nosuch-user/", headers)

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

const timesheet197604 = `1976 April
----------
14  1 Thu 8
    2 Fri 8
    3 Sat
    4 Sun
15  5 Mon 8
    6 Tue 8
    7 Wed 8
    8 Thu 8
    9 Fri 8
   10 Sat
   11 Sun
16 12 Mon 8
   13 Tue 8
   14 Wed 8
   15 Thu 8
   16 Fri 8
   17 Sat
   18 Sun
17 19 Mon 8
   20 Tue 8
   21 Wed 8
   22 Thu 8
   23 Fri 8
   24 Sat
   25 Sun
18 26 Mon 8
   27 Tue 8
   28 Wed 8
   29 Thu 8
   30 Fri 8
`
