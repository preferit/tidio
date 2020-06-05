package permission

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_rules(t *testing.T) {
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	ok(ToRead(0, 0, &thing{0, 0, UserRWX}))
	ok(ToRead(0, 0, &thing{0, 0, GroupRWX}))
	ok(ToRead(1, 1, &thing{0, 0, OtherR}))

	bad(ToRead(0, 0, &thing{0, 0, NoMode}))
	bad(ToRead(1, 0, &thing{0, 0, UserRWX}))
}
