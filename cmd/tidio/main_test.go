package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestCommand_run_default(t *testing.T) {
	cmd := NewCommand()
	cmd.ListenAndServe = func(string, http.Handler) error { return nil }

	ok, bad := asserter.NewErrors(t)

	ok(cmd.run("name"))
	os.RemoveAll("system.state")

	ok(cmd.run("name", "-state", "somefile"))
	ok(cmd.run("name", "-state", "somefile")) // should reload
	os.RemoveAll("somefile")

	bad(cmd.run("name", "-no-such"))
	bad(cmd.run("name", "-state", "/no-such"))

	w, err := ioutil.TempFile("", "tidio")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(w, "jibberish")
	w.Close()
	defer os.RemoveAll(w.Name())
	bad(cmd.run("name", "-state", w.Name()))

}
