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
		account      = Account{}
	)
	defer cleanup()
	service.SetDataDir(dir)
	john.Key = key
	service.AddAccount(john)

	bad(service.Load())
	ok(service.Save())
	ok(service.Load())

	ok(service.AccountByKey(&account, key))
	bad(service.AccountByKey(&account, ""))
	bad(service.AccountByKey(&account, "not there"))
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
