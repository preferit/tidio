package tidio

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_service(t *testing.T) {
	var (
		assert       = asserter.New(t)
		ok, bad      = assert().Errors()
		service      = newTestService(t)
		dir, cleanup = newTempDir(t)
	)
	defer cleanup()

	bad(service.LoadState(path.Join(dir, "data.gob")))
	ok(service.SaveState())

	if _, ok := service.RoleByKey("KEY"); !ok {
		t.Error("KEY is in apikeys")
	}
	if _, ok := service.RoleByKey(""); ok {
		t.Error("empty key ok")
	}
	if _, ok := service.RoleByKey("not there"); ok {
		t.Error("wrong key ok")
	}
}

func newTestService(t *testing.T) *Service {
	apikeys := APIKeys{
		"KEY": NewAccount("john", "admin"),
	}
	service := NewService(apikeys)
	return service
}

func newTempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "tidioservice")
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func Test_service_options(t *testing.T) {
	defer catchPanic(t)
	NewService(1)
}

func catchPanic(t *testing.T) {
	e := recover()
	if e == nil {
		t.Helper()
		t.Error("didn't panic")
	}
}
