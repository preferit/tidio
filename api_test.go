package tidio

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/asserter"
	"github.com/gregoryv/fox"
	"github.com/gregoryv/go-timesheet"
)

var (
	withJohnAccount = InitialAccount{"john", "secret"}
	asJohn          = NewCredentials("john", "secret")
)

func TestClient_CreateTimesheet_asJohn(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv)
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn, log)
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
		resp = api.CreateTimesheet(path, body).MustSend()
	)
	if resp.StatusCode != 201 {
		t.Error(resp.Status)
	}
}

func TestClient_CreateTimesheet_asAnonymous(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv)
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, log)
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
		resp = api.CreateTimesheet(path, body).MustSend()
	)
	if resp.StatusCode != 401 {
		t.Error(resp.Status)
	}
}

func TestClient_ReadTimesheet_asJohn(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv)
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn)
		path = "/api/timesheets/john/202001.timesheet"
		body = timesheet.Render(2020, 1, 8)
		_    = api.CreateTimesheet(path, body).MustSend()
		resp = api.ReadTimesheet(path).MustSend()
	)
	if resp.StatusCode != 200 {
		t.Error(resp.Status)
	}
}

func TestClient_ReadTimesheet_noSuchResource(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv)
	)
	defer ts.Close()
	var (
		api  = NewAPI(ts.URL, asJohn)
		resp = api.ReadTimesheet("/api/timesheets/john/nosuch").MustSend()
	)
	if resp.StatusCode != 404 {
		t.Error(resp.Status)
	}
}

func TestClient_ReadTimesheet_asAnonymous(t *testing.T) {
	var (
		log = fox.Logging{t}
		srv = NewService(log, withJohnAccount)
		ts  = httptest.NewServer(srv)
	)
	defer ts.Close()
	var (
		api       = NewAPI(ts.URL, asJohn)
		path      = "/api/timesheets/john/202001.timesheet"
		body      = timesheet.Render(2020, 1, 8)
		_         = api.CreateTimesheet(path, body).MustSend()
		anonymous = NewCredentials("", "")
	)
	ant.MustConfigure(api, anonymous)
	resp := api.ReadTimesheet(path).MustSend()
	if resp.StatusCode == 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		t.Error("should fail\n", buf.String())
	}
}

func Test_defaults(t *testing.T) {
	srv := NewService()
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(200, "GET", "/")
	exp.StatusCode(200, "GET", "/api")
	exp.StatusCode(405, "X", "/api")
	exp.StatusCode(http.StatusUnauthorized, "GET", "/api/timesheets")
	exp.BodyIs(`{"resources":[{"name": "timesheets"}]}`, "GET", "/api")
}
