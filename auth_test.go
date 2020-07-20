package tidio

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestAuth_Token(t *testing.T) {
	auth := NewAuth()
	ok, bad := asserter.NewErrors(t)
	ok(auth.Parse("Bearer xxx"))
	bad(auth.Parse(""))
}
