package tidio

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/rs"
)

func TestService_ServeHTTP_POST_timesheet_asJohn(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	srv.AddAccount("john", "secret")

	w := httptest.NewRecorder()
	validTimesheet := strings.NewReader("todo")
	r, _ := http.NewRequest("POST", "/api/timesheets/john/202001.timesheet", validTimesheet)
	r.SetBasicAuth("john", "secret")
	srv.ServeHTTP(w, r)
	got := w.Result()

	if got.StatusCode != 201 {
		body, _ := ioutil.ReadAll(got.Body)
		t.Error(got.Status, string(body))
	}
}

func TestService_ServeHTTP_POST_timesheets_missing_body_asJohn(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	srv.AddAccount("john", "secret")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/timesheets/john", nil)
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

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/timesheets/john", nil)
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

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/timesheets/john", nil)
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

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/timesheets/john", nil)
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
