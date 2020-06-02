package main

import (
	"net/http"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/preferit/tidio"
)

func Test_router(t *testing.T) {
	assert := asserter.New(t)
	apikeys := map[string]string{
		"KEY": "john",
	}
	headers := http.Header{}
	store := tidio.NewStore("")

	exp := assert().ResponseFrom(NewRouter(apikeys, store))
	exp.StatusCode(200, "GET", "/api", nil)
	exp.Contains("revision", "GET", "/api")
	exp.Contains("version", "GET", "/api")
	exp.Contains("resources", "GET", "/api")

	exp.StatusCode(401, "GET", "/api/timesheets/")
	headers = http.Header{}
	headers.Set("Authorization", "NO SUCH KEY")
	exp.StatusCode(401, "GET", "/api/timesheets/", headers)
	exp.StatusCode(401, "POST", "/api/timesheets/")

	// authenticated
	headers = http.Header{}
	headers.Set("Authorization", "KEY")
	exp.StatusCode(200, "GET", "/api/timesheets/", headers)
}
