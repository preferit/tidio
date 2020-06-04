package tidio

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gregoryv/asserter"
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
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	ok(store.WriteFile("john", filename, aFile(content)))
	os.Chmod(path.Join(store.dir, filename), 0400)
	bad(store.WriteFile("john", filename, aFile(content)), "read only")

	var sheet Timesheet
	ok(store.OpenFile(&sheet, filename))
	body, err := ioutil.ReadAll(&sheet)
	ok(err)
	sheet.Close()
	assert().Equals(string(body), content)

	bad(store.OpenFile(&Timesheet{}, "no such file"))
}

func Test_store_Glob(t *testing.T) {
	store, cleanup := newTempStore(t)
	defer cleanup()
	store.Init()
	store.WriteFile("john", "john/file1", aFile(""))
	store.WriteFile("eva", "john/file2", aFile(""))
	files := store.Glob("john", "*")
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
