package tidio

import (
	"io/ioutil"
	"os"
	"testing"
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

func Test_store_fileio(t *testing.T) {
	store, cleanup := newTempStore(t)
	defer cleanup()
	store.Init()
	filename := "a/b/something.x"

	if err := store.WriteFile(filename, []byte(".."), 0644); err != nil {
		t.Fatal(err)
	}
	if err := store.ReadFile(ioutil.Discard, filename); err != nil {
		t.Error(err)
	}

	t.Run("no such file", func(t *testing.T) {
		if err := store.ReadFile(ioutil.Discard, "no such file"); err == nil {
			t.Error("did not fail")
		}
	})
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
