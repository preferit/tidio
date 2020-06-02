package tidio

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_store(t *testing.T) {
	store, cleanup := newTempStore(t)
	defer cleanup()

	if store.IsInitiated() {
		t.Error("new store should be uninitiated")
	}
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	if !store.IsInitiated() {
		t.Fail()
	}
}

func Test_store_writefile(t *testing.T) {
	store, cleanup := newTempStore(t)
	defer cleanup()
	store.Init()
	if err := store.WriteFile("a/b/something.x", []byte(".."), 0644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
}

func newTempStore(t *testing.T) (*Store, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "tidiostore")
	if err != nil {
		t.Fatal(err)
	}
	store := NewStore(dir)
	store.Logger = t
	return store, func() { os.RemoveAll(dir) }
}
