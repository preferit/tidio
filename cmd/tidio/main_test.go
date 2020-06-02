package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func Test_cli(t *testing.T) {
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

	t.Run("run", func(t *testing.T) {
		c := &cli{keysfile: fh.Name()}
		if err := c.run(); err == nil {
			t.Fail()
		}
	})

	t.Run("start", func(t *testing.T) {
		var started bool
		c := &cli{
			bind:     ":8080",
			keysfile: fh.Name(),
			starter:  func(string, http.Handler) error { started = true; return nil },
		}
		if err := c.run(); err != nil {
			t.Fatal(err)
		}
		if !started {
			t.Error("starter not called")
		}
	})
}
