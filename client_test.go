package tidio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/asserter"
	"github.com/gregoryv/go-timesheet"
)

func TestClient_CreateTimesheet_asJohn(t *testing.T) {
	client := NewClient(Logging{t}, ErrorHandling(t.Fatal))
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	cred := NewCredentials("john", "secret")
	api := NewAPI(ts.URL, cred)
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	r, _ := api.CreateTimesheet(path, body)

	resp, _ := client.Send(r)
	if resp.StatusCode != 201 {
		t.Error(resp.Status)
	}
}

func TestClient_CreateTimesheet_asAnonymous(t *testing.T) {
	client := NewClient(Logging{t}, ErrorHandling(t.Fatal))
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	api := NewAPI(ts.URL)
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	r, _ := api.CreateTimesheet(path, body)

	resp, _ := client.Send(r)
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
	client := NewClient(Logging{t}, ErrorHandling(t.Fatal))
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	cred := NewCredentials("john", "secret")
	api := NewAPI(ts.URL, cred)
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	r, _ := api.CreateTimesheet(path, body)

	client.Send(r)

	r, _ = api.ReadTimesheet(path)
	client.Send(r)
}

func TestClient_ReadTimesheet_asAnonymous(t *testing.T) {
	client := NewClient(Logging{t}, ErrorHandling(t.Fatal))
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	cred := NewCredentials("john", "secret")
	api := NewAPI(ts.URL, cred)
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	r, _ := api.CreateTimesheet(path, body)
	client.Send(r)

	anonymous := NewCredentials("", "")
	ant.MustConfigure(api, anonymous)

	r, _ = api.ReadTimesheet(path)
	resp, _ := client.Send(r)

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

func TestClient_with_bad_setting(t *testing.T) {
	defer catchPanic(t)
	NewClient(ant.SetFunc(func(interface{}) error {
		return fmt.Errorf("bad client setting")
	}))
	t.Error("should panic")
}
