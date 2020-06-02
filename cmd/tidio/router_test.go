package main

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_router(t *testing.T) {
	assert := asserter.New(t)
	exp := assert().ResponseFrom(NewRouter())
	exp.StatusCode(200, "GET", "/api", nil)
	exp.Contains("revision", "GET", "/api", nil)
	exp.Contains("version", "GET", "/api", nil)
}