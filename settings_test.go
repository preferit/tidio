package tidio

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestSettings(t *testing.T) {
	bad := asserter.Wrap(t).Bad
	bad(InitialAccount{}.Set(nil))
}
