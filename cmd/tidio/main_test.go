package main

import (
	"net/http"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_cli(t *testing.T) {
	var (
		assert  = asserter.New(t)
		ok, bad = assert().Errors()
	)
	bad(run(&cli{Logger: t}))

	ok(run(&cli{
		bind: ":8080",
		ListenAndServe: func(string, http.Handler) error {
			return nil
		},
		Logger: t,
	}))
}
