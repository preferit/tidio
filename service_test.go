package tidio

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/ex"
)

func TestService_NewAccount(t *testing.T) {
	var (
		service = NewService()
		assert  = asserter.New(t)
		ok, bad = assert().Errors()
		account Account
	)
	ok(service.NewAccount(&account, "john"))
	bad(service.NewAccount(&account, ""))
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

	assert().Contains(buf.String(), "resources")
	assert().Contains(buf.String(), "Path")
	assert().Contains(buf.String(), "Entity")
	assert().Contains(buf.String(), "Mode")
	assert().Contains(buf.String(), "UID")
	assert().Contains(buf.String(), "GID")
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
