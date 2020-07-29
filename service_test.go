package tidio

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/go-timesheet"
	"github.com/gregoryv/rs"
)

func TestService_ServeHTTP_GET_timesheet_asJohn(t *testing.T) {
	srv := NewService()
	srv.AddAccount("john", "secret")

	var sheet bytes.Buffer
	timesheet.Render(&sheet, 2020, 1, 8)
	exp := sheet.Bytes()
	path := "/api/timesheets/john/202001.timesheet"
	w, r := wr(t, "POST", path, bytes.NewReader(exp))
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)

	w, r = wr(t, "GET", path, nil)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(body, exp) {
		t.Error(string(body))
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

func TestService_ServeHTTP_POST_timesheet_asJohn(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	srv.AddAccount("john", "secret")

	var sheet bytes.Buffer
	timesheet.Render(&sheet, 2020, 1, 8)
	w, r := wr(t, "POST", "/api/timesheets/john/202001.timesheet", &sheet)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	resp := w.Result()

	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Error(resp.Status, string(body))
	}
}

func TestService_ServeHTTP_POST_timesheets_missing_body_asJohn(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	srv.AddAccount("john", "secret")

	w, r := wr(t, "POST", "/api/timesheets/john", nil)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 400 {
		t.Error(got.Status)
	}
}

func TestService_ServeHTTP_GET_timesheets_authenticated(t *testing.T) {
	srv := NewService()
	srv.AddAccount("john", "secret")

	w, r := wr(t, "GET", "/api/timesheets/john", nil)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 200 {
		t.Error(got.Status)
	}
}

func TestService_ServeHTTP_GET_timesheets_badheader(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	srv.AddAccount("john", "secret")

	w, r := wr(t, "GET", "/api/timesheets/john", nil)
	r.Header.Set("Authorization", "Basic jibberish")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 401 {
		t.Error(got.Status)
	}
}

func TestService_ServeHTTP_GET_timesheets_autherror(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	srv.AddAccount("john", "secret")

	w, r := wr(t, "GET", "/api/timesheets/john", nil)
	r.SetBasicAuth("john", "wrong")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 401 {
		t.Error(got.Status)
	}
}

func TestService_AddAccount(t *testing.T) {
	srv := NewService()
	ok, bad := asserter.NewErrors(t)
	ok(srv.AddAccount("john", "secret"))
	bad(srv.AddAccount("john", "secret"))
	bad(srv.AddAccount("root", "secret"))
	_, err := rs.Root.Use(srv.sys).Stat("/api/timesheets/john")
	ok(err)
}

func TestService_ServeHTTP(t *testing.T) {
	srv := NewService()
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(200, "GET", "/api")
	exp.StatusCode(405, "X", "/api")
	exp.BodyIs(`{"resources":[{"name": "timesheets"}]}`, "GET", "/api")
}

func TestService_anonymousAccess(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(http.StatusUnauthorized, "GET", "/api/timesheets")
}
