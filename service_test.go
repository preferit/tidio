package tidio

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_service(t *testing.T) {
	var (
		assert       = asserter.New(t)
		ok, bad      = assert().Errors()
		service      = NewService()
		dir, cleanup = newTempDir(t)
		john         = NewAccount("john")
		key          = "KEY"
	)
	defer cleanup()
	service.SetDataDir(dir)
	john.Key = key
	service.AddAccount(john)

	bad(service.Load())
	ok(service.Save())
	ok(service.Load())

	if _, ok := service.AccountByKey(key); !ok {
		t.Error("KEY is in apikeys")
	}
	if _, ok := service.AccountByKey(""); ok {
		t.Error("empty key ok")
	}
	if _, ok := service.AccountByKey("not there"); ok {
		t.Error("wrong key ok")
	}
}

func newTempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "tidioservice")
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func catchPanic(t *testing.T) {
	e := recover()
	if e == nil {
		t.Helper()
		t.Error("didn't panic")
	}
}
