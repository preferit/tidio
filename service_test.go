package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
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

func TestService_options(t *testing.T) {
	defer catchPanic(t)
	NewService(
		ServiceOptFunc(func(*Service) error {
			return fmt.Errorf("option failed")
		}),
	)
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
