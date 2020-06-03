package tidio

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
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
	content := "body"

	if err := store.WriteFile(filename, aFile(content)); err != nil {
		t.Fatal(err)
	}
	os.Chmod(path.Join(store.dir, filename), 0400)
	if err := store.WriteFile(filename, aFile(content)); err == nil {
		t.Error("wrote read only file")
	}
	var buf bytes.Buffer
	if err := store.ReadFile(&buf, filename); err != nil {
		t.Error(err)
	}
	got := buf.String()
	if got != content {
		t.Error("wrong content:", got)
	}
	t.Run("no such file", func(t *testing.T) {
		if err := store.ReadFile(ioutil.Discard, "no such file"); err == nil {
			t.Error("did not fail")
		}
	})
}

func Test_store_Glob(t *testing.T) {
	store, cleanup := newTempStore(t)
	defer cleanup()
	store.Init()
	store.WriteFile("john/file1", aFile(""))
	store.WriteFile("john/file2", aFile(""))
	files, err := store.Glob("john", "*")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Error("expected 2 files:", files)
	}
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
