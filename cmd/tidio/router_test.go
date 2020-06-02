package main

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_router(t *testing.T) {
	assert := asserter.New(t)
	apikeys := make(map[string]string)
	exp := assert().ResponseFrom(NewRouter(apikeys))
	exp.StatusCode(200, "GET", "/api", nil)
	exp.Contains("revision", "GET", "/api")
	exp.Contains("version", "GET", "/api")
	exp.Contains("resources", "GET", "/api")

	exp.StatusCode(401, "GET", "/api/timesheets/")
	exp.StatusCode(401, "POST", "/api/timesheets/")
}
