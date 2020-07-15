package tidio

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestService_ServeHTTP(t *testing.T) {
	service := NewService()
	assert := asserter.New(t)
	exp := assert().ResponseFrom(service)
	exp.StatusCode(200, "GET", "/api")
}
