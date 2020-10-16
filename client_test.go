package tidio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/go-timesheet"
)

func TestClient_CreateTimesheet_asJohn(t *testing.T) {
	srv := NewService(Logging{t}, InitialAccount{"john", "secret"})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	api := API{host: ts.URL}
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	r, _ := api.CreateTimesheet(path, body)
	client := NewClient(
		Logging{t},
	)
	cred := Credentials{account: "john", secret: "secret"}
	_, err := client.Send(r, &cred)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadTimesheet_asJohn(t *testing.T) {
	srv := NewService(
		Logging{t}, InitialAccount{"john", "secret"},
	)
	ts := httptest.NewServer(srv)
	defer ts.Close()
	api := NewAPI(ts.URL)
	client := NewClient(
		Logging{t},
		ErrorHandling(t.Fatal),
	)
	cred := &Credentials{account: "john", secret: "secret"}

	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	r, _ := api.CreateTimesheet(path, body)
	client.Send(r, cred)

	r, _ = api.ReadTimesheet(path)
	client.Send(r, cred)
}

func TestClient_ReadTimesheet_asAnonymous(t *testing.T) {
	srv := NewService(
		Logging{t}, InitialAccount{"john", "secret"},
	)
	ts := httptest.NewServer(srv)
	defer ts.Close()
	client := NewClient(
		Logging{t},
		ErrorHandling(t.Fatal),
	)
	api := NewAPI(ts.URL)
	cred := &Credentials{account: "john", secret: "secret"}
	path := "/api/timesheets/john/202001.timesheet"
	body := timesheet.Render(2020, 1, 8)
	r, _ := api.CreateTimesheet(path, body)
	client.Send(r, cred)

	r, _ = api.ReadTimesheet(path)
	resp, _ := client.Send(r, nil) // anonymous, nil credentials
	if resp.StatusCode == 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		t.Error("should fail\n", buf.String())
	}
}

func wr(t *testing.T, method, url string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()
	w := httptest.NewRecorder()
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return w, r
}

func Test_POST_timesheets_missing_body_asJohn(t *testing.T) {
	srv, output := newTestService()

	w, r := wr(t, "POST", "/api/timesheets/john", nil)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 400 {
		t.Error(got.Status, output)
	}
}

func Test_GET_timesheets_badheader(t *testing.T) {
	srv, output := newTestService()

	w, r := wr(t, "GET", "/api/timesheets/john", nil)
	r.Header.Set("Authorization", "Basic jibberish")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 401 {
		t.Error(got.Status, output)
	}
}

func Test_GET_timesheets_autherror(t *testing.T) {
	srv, output := newTestService()

	w, r := wr(t, "GET", "/api/timesheets/john", nil)
	r.SetBasicAuth("john", "wrong")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 401 {
		t.Error(got.Status, output)
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
	NewClient(SetFunc(func(interface{}) error {
		return fmt.Errorf("bad client setting")
	}))
	t.Error("should panic")
}

func newTestService() (*Service, *BufferedLogger) {
	srv := NewService()
	buflog := Buflog(srv)
	srv.AddAccount("john", "secret")
	return srv, buflog
}
