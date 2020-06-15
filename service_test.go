package tidio

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/ex"
)

func TestService_AddUser(t *testing.T) {
	var (
		service = NewService()
		assert  = asserter.New(t)
		ok, bad = assert().Errors()
	)
	ok(service.AddUser(&Account{Username: "something"}))
	bad(service.AddUser(&Account{}))
}

func TestService_WriteTo(t *testing.T) {
	var (
		assert  = asserter.New(t)
		service = NewService()
		buf     bytes.Buffer
		nice    = ex.NewJsonWriter()
	)
	nice.Out = &buf
	service.WriteTo(nice)
	got := buf.String()
	assert().Contains(got, "resources")
	assert().Contains(got, "Path")
	assert().Contains(got, "Entity")
	assert().Contains(got, "Mode")
	assert().Contains(got, "UID")
	assert().Contains(got, "GID")
	assert().Contains(got, "Username")
	assert().Contains(got, `"root"`)
}

// ----------------------------------------

func TestService_ServeHTTP(t *testing.T) {
	var (
		assert  = asserter.New(t)
		service = NewService()
		exp     = assert().ResponseFrom(service)
	)
	exp.StatusCode(200, "GET", "/api")
}
