package main

import (
	"net/http"
	"testing"
)

func Test_cli(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		c := &cli{}
		if err := c.run(); err == nil {
			t.Fail()
		}
	})

	t.Run("start", func(t *testing.T) {
		var started bool
		c := &cli{
			bind:    ":8080",
			starter: func(string, http.Handler) error { started = true; return nil },
		}
		if err := c.run(); err != nil {
			t.Fatal(err)
		}
		if !started {
			t.Error("starter not called")
		}
	})

}
