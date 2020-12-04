package tidio

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/ant"
	"github.com/gregoryv/asserter"
	"github.com/gregoryv/rs"
)

func TestService_AddAccount(t *testing.T) {
	srv := NewService()
	ok := asserter.Wrap(t).Ok
	ok(srv.AddAccount("john", "secret"))

	bad := asserter.Wrap(t).Bad
	bad(srv.AddAccount("john", "secret"))
	bad(srv.AddAccount("root", "secret"))
	bad(srv.AddAccount("eva", ""))

	xok := asserter.Wrap(t).MixOk
	john := rs.NewAccount("john", 2)
	asJohn := john.Use(srv.sys)
	xok(asJohn.Stat("/api/timesheets/john"))
}

func TestService_RestoreState_missing_file(t *testing.T) {
	srv := NewService()
	ok := asserter.Wrap(t).Ok
	ok(srv.RestoreState())
}

func TestNewService_panics_on_bad_settings(t *testing.T) {
	defer catchPanic(t)
	NewService(
		ant.SetFunc(func(interface{}) error {
			return fmt.Errorf("option failed")
		}),
	)
	t.Error("should panic")
}

func newTestService() (*Service, *BufferedLogger) {
	srv := NewService()
	buflog := Buflog(srv)
	srv.AddAccount("john", "secret")
	return srv, buflog
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

func catchPanic(t *testing.T) {
	e := recover()
	if e == nil {
		t.Error("expected panic")
	}
}

type brokenStorage struct {
	called bool
}

// Create
func (me *brokenStorage) Create() (io.WriteCloser, error) {
	me.called = true
	return nil, fmt.Errorf("brokenStorage")
}
