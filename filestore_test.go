package tidio

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_store(t *testing.T) {
	dir, err := ioutil.TempDir("", "tidiostore")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store := NewStore(dir)
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
