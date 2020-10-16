package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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
	asRoot := rs.Root.Use(srv.sys)
	xok(asRoot.Stat("/api/timesheets/john"))
}

func TestService_RestoreState_missing_file(t *testing.T) {
	srv := NewService()
	bad := asserter.Wrap(t).Bad
	bad(srv.RestoreState("no-such-file"))
}

func TestService_RestoreState_ok_file(t *testing.T) {
	srv := NewService()
	tmp, _ := ioutil.TempFile("", "restorestate")
	srv.sys.Export(tmp)
	tmp.Close()
	ok := asserter.Wrap(t).Ok
	ok(srv.RestoreState(tmp.Name()))
	os.RemoveAll(tmp.Name())
}

func TestService_AutoPersist(t *testing.T) {
	srv := NewService()
	buflog := Buflog(srv)

	tmp, _ := ioutil.TempFile("", "restorestate")
	tmp.Close()
	defer os.RemoveAll(tmp.Name())

	dest := NewFileStorage(tmp.Name())
	srv.AutoPersist(dest, time.Millisecond)

	// make a change
	asRoot := rs.Root.Use(srv.sys)
	asRoot.Exec("/bin/mkdir /tmp/x")

	time.Sleep(10 * time.Millisecond)
	got, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err, "\n", buflog.String())
	}

	assert := asserter.New(t)
	assert(len(got) != 0).Error("empty state", "\n", buflog.String())
}

func TestService_AutoPersist_create_file_fails(t *testing.T) {
	srv := NewService()
	buflog := Buflog(srv)

	dest := &brokenStorage{}
	srv.AutoPersist(dest, time.Millisecond)

	// make a change
	asRoot := rs.Root.Use(srv.sys)
	asRoot.Exec("/bin/mkdir /tmp/x")

	time.Sleep(10 * time.Millisecond)
	assert := asserter.New(t)
	assert(dest.called).Error("state persisted", "\n", buflog.String())
}

func TestNewService_panics_on_bad_settings(t *testing.T) {
	defer catchPanic(t)
	NewService(
		SetFunc(func(interface{}) error {
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
