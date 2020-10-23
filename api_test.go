package tidio

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/asserter"
	"github.com/gregoryv/go-timesheet"
)

func TestClient_CreateTimesheet_asJohn(t *testing.T) {
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	cred := NewCredentials("john", "secret")
	api := NewAPI(ts.URL, cred, Logging{t})
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	resp, _ := api.CreateTimesheet(path, body).Send()

	if resp.StatusCode != 201 {
		t.Error(resp.Status)
	}
}

func TestClient_CreateTimesheet_asAnonymous(t *testing.T) {
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	api := NewAPI(ts.URL)
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	resp, _ := api.CreateTimesheet(path, body).Send()

	if resp.StatusCode != 401 {
		t.Error(resp.Status)
	}
}

func dump(r io.Reader) string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	io.Copy(&buf, r)
	return buf.String()
}

func TestClient_ReadTimesheet_asJohn(t *testing.T) {
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	cred := NewCredentials("john", "secret")
	api := NewAPI(ts.URL, cred)
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	api.CreateTimesheet(path, body).Send()

	resp, _ := api.ReadTimesheet(path).Send()
	if resp.StatusCode != 200 {
		t.Error(resp.Status)
	}
}

func TestClient_ReadTimesheet_noSuchResource(t *testing.T) {
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	cred := NewCredentials("john", "secret")
	api := NewAPI(ts.URL, cred)

	resp, _ := api.ReadTimesheet("/api/timesheets/john/nosuch").Send()
	if resp.StatusCode != 404 {
		t.Error(resp.Status)
	}
}

func TestClient_ReadTimesheet_asAnonymous(t *testing.T) {
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	cred := NewCredentials("john", "secret")
	api := NewAPI(ts.URL, cred)
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	api.CreateTimesheet(path, body).Send()

	anonymous := NewCredentials("", "")
	ant.MustConfigure(api, anonymous)

	resp, _ := api.ReadTimesheet(path).Send()

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
