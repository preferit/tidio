package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestCommand_run_default(t *testing.T) {
	cmd := NewCommand()
	cmd.Logger = t
	cmd.ListenAndServe = func(string, http.Handler) error { return nil }

	ok := asserter.Wrap(t).Ok
	ok(cmd.run("name", "-h"))

	ok(cmd.run("name"))
	defer os.RemoveAll("system.state")

	name, cleanup := writeTempFile(t, "")
	defer cleanup()
	os.RemoveAll(name)
	t.Log(name)
	bad := asserter.Wrap(t).Bad
	ok(cmd.run("name", "-state", name))
	ok(cmd.run("name", "-state", name)) // should reload
	return

	//	bad := asserter.Wrap(t).Bad
	bad(cmd.run("name", "-no-such"))
	bad(cmd.run("name", "-state", "/no-such"))

	// badly formatted state file
	name, cleanup = writeTempFile(t, "jibberish")
	defer cleanup()
	bad(cmd.run("name", "-state", name))
}

func writeTempFile(t *testing.T, content interface{}) (name string, cleanup func()) {
	t.Helper()
	w, err := ioutil.TempFile("", "tidio")
	if err != nil {
		t.Fatal(err)
	}

	switch content := content.(type) {
	case io.Reader:
		io.Copy(w, content)
	default:
		fmt.Fprintf(w, "%v", content)
	}
	w.Close()
	return w.Name(), func() {
		os.RemoveAll(w.Name())
	}
}
