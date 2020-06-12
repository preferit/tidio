package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/preferit/tidio"
)

func Test_cli(t *testing.T) {
	badFile, err := ioutil.TempFile("", "apikeys")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(badFile, "{...") // bad json
	badFile.Close()
	defer os.RemoveAll(badFile.Name())

	bad := func(msg string, c *cli) {
		t.Helper()
		if err := c.run(); err == nil {
			t.Error(msg)
		}
	}
	bad("bind not set", &cli{})
	bad("keysfile not found", &cli{
		bind:     "1",
		keysfile: "NO SUCH FILE",
	})
	bad("bad formata keysfile", &cli{
		bind:     "1",
		keysfile: badFile.Name(),
	})

	// create correct key file
	fh, err := ioutil.TempFile("", "apikeys")
	if err != nil {
		t.Fatal(err)
	}
	accounts := tidio.AccountsMap{}.New()
	accounts.AddAccount("KEY", tidio.NewAccount("john", "admin"))
	accounts.SaveAccounts(fh)
	defer os.RemoveAll(fh.Name())

	// setup store
	storeDir, err := ioutil.TempDir("", "tidio")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(storeDir)
	ok := func(c *cli) {
		t.Helper()
		if err := c.run(); err != nil {
			t.Error(err)
		}
	}
	ok(&cli{
		storeDir: storeDir,
		bind:     ":80",
		keysfile: fh.Name(),
		starter:  func(string, http.Handler) error { return nil },
	})
}

func newTempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "tidiocmd")
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}
