package tidio

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestService(t *testing.T) {
	var (
		service = NewService()
		assert  = asserter.New(t)
		ok, bad = assert().Errors()
		account Account
	)
	ok(service.NewAccount(&account, "john"))
	bad(service.NewAccount(&account, ""))
}

func TestService_ServeHTTP(t *testing.T) {
	var (
		assert  = asserter.New(t)
		service = NewService()
		exp     = assert().ResponseFrom(service)
	)

	exp.StatusCode(200, "GET", "/api")
}

func TestService_WriteTo(t *testing.T) {
	var (
		assert  = asserter.New(t)
		service = NewService()
		buf     bytes.Buffer
	)
	service.WriteTo(&buf)
	assert().Contains(buf.String(), "resources")
}
