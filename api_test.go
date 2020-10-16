package tidio

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/go-timesheet"
)

func TestAPI_CreateTimesheet_asJohn(t *testing.T) {
	srv := NewService(
		LoggerOption{t},
	)
	srv.AddAccount("john", "secret")

	ts := httptest.NewServer(srv)
	defer ts.Close()

	v1 := NewAPIv1(ts.URL, Credentials{account: "john", secret: "secret"})
	ok := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	ok(v1.CreateTimesheet(
		"/api/timesheets/john/202001.timesheet",
		timesheet.Render(2020, 1, 8),
	))

	bad := func(err error) {
		t.Helper()
		if err == nil {
			t.Fatal("should fail")
		}
	}
	bad(v1.CreateTimesheet(
		"/NOTOK/timesheets/john/202001.timesheet",
		timesheet.Render(2020, 1, 8),
	))
	_ = bad

}

func Test_GET_timesheet_asJohn(t *testing.T) {
	var (
		srv     = NewService()
		_       = srv.AddAccount("john", "secret")
		exp     = asserter.Wrap(t).ResponseFrom(srv)
		headers = basicAuthHeaders("john", "secret")
		sheet   bytes.Buffer
		path    = "/api/timesheets/john/202001.timesheet"
	)
	timesheet.RenderTo(&sheet, 2020, 1, 8)
	sheetStr := sheet.String()

	exp.StatusCode(201, "POST", path, strings.NewReader(sheetStr), headers)
	exp.StatusCode(200, "GET", path, headers)
	exp.BodyIs(sheetStr, "GET", path, headers)
	// FIXME
	//exp.StatusCode(403, "POST", "/api/timesheets/202001.timesheet", &sheet, headers)
	exp.Contains("denied", "POST", "/api/timesheets/202001.timesheet",
		strings.NewReader(sheetStr), headers)

	exp.StatusCode(404, "GET", "/api/no-such-path", headers)
	// FIXME
	//exp.StatusCode(403, "GET", "/etc/accounts/john", headers)
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

func Test_POST_timesheet_asJohn(t *testing.T) {
	srv, output := newTestService()

	body := timesheet.Render(2020, 1, 8)
	w, r := wr(t, "POST", "/api/timesheets/john/202001.timesheet", body)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	resp := w.Result()

	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Error(resp.Status, string(body), output)
	}
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

func Test_GET_timesheets_authenticated(t *testing.T) {
	srv, output := newTestService()

	w, r := wr(t, "GET", "/api/timesheets/john", nil)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 200 {
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

func newTestService() (*Service, *BufferedLogger) {
	srv := NewService()
	buflog := Buflog(srv)
	srv.AddAccount("john", "secret")
	return srv, buflog
}
