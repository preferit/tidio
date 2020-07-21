package tidio

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestKey_Check(t *testing.T) {
	key := NewKey("secret", "/etc/accounts/john.acc")
	ok, bad := asserter.NewErrors(t)
	ok(key.Check("secret"))
	bad(key.Check(""))
}
