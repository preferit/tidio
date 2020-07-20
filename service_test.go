package tidio

import (
	"net/http"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestService_ServeHTTP(t *testing.T) {
	srv := NewService()
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(200, "GET", "/api")
	exp.StatusCode(405, "X", "/api")
	exp.BodyIs(`{"resources":[{"name": "timesheets"}]}`, "GET", "/api")
}

func TestService_anonymousAccess(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(http.StatusUnauthorized, "GET", "/api/timesheets")
}
