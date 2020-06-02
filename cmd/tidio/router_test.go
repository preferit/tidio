package main

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_router(t *testing.T) {
	assert := asserter.New(t)
	exp := assert().ResponseFrom(NewRouter())
	exp.StatusCode(200, "GET", "/", nil)
}
