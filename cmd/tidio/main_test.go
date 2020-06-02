package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
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
	apikeys := map[string]string{
		"john": "KEY",
	}
	json.NewEncoder(fh).Encode(apikeys)
	fh.Close()
	defer os.RemoveAll(fh.Name())

	ok := func(c *cli) {
		t.Helper()

		if err := c.run(); err != nil {
			t.Error(err)
		}
	}
	ok(&cli{
		bind:     ":80",
		keysfile: fh.Name(),
		starter:  func(string, http.Handler) error { return nil },
	})
}
