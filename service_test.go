package tidio

import (
	"net/http"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/rs"
)

func TestService_AddAccount(t *testing.T) {
	srv := NewService()
	ok, bad := asserter.NewErrors(t)
	ok(srv.AddAccount("john", "secret"))
	bad(srv.AddAccount("john", "secret"))
	bad(srv.AddAccount("root", "secret"))
	_, err := rs.Root.Use(srv.sys).Stat("/api/timesheets/john")
	ok(err)
}

func TestService_ServeHTTP(t *testing.T) {
	srv := NewService()
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(200, "GET", "/api")
	exp.StatusCode(405, "X", "/api")
	exp.BodyIs(`{"resources":[{"name": "timesheets"}]}`, "GET", "/api")
}

func TestService_ServeHTTP_authenticated(t *testing.T) {
	srv := NewService()
	john := &BasicAuth{AccountName: "john", Secret: "secret"}
	srv.AddAccount("john", "secret")
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(200, "GET", "/api/timesheets/john", http.Header{
		"Authorization": []string{"Basic " + john.Token()},
	})
}

func TestService_anonymousAccess(t *testing.T) {
	srv := NewService()
	srv.SetLogger(t)
	assert := asserter.New(t)
	exp := assert().ResponseFrom(srv)
	exp.StatusCode(http.StatusUnauthorized, "GET", "/api/timesheets")
}
