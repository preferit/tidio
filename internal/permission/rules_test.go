package permission

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_rules(t *testing.T) {
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	ok(ToRead(&thing{0, 0, UserR}, &mockAccount{0, []int{0}}))
	ok(ToRead(&thing{0, 0, GroupR}, &mockAccount{0, []int{0}}))
	ok(ToRead(&thing{0, 0, OtherR}, &mockAccount{1, []int{1}}))
	bad(ToRead(&thing{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToRead(&thing{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToWrite(&thing{0, 0, UserW}, &mockAccount{0, []int{0}}))
	ok(ToWrite(&thing{0, 0, GroupW}, &mockAccount{0, []int{0}}))
	ok(ToWrite(&thing{0, 0, OtherW}, &mockAccount{1, []int{1}}))
	bad(ToWrite(&thing{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToWrite(&thing{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToExec(0, 0, &thing{0, 0, UserX}))
	ok(ToExec(0, 0, &thing{0, 0, GroupX}))
	ok(ToExec(1, 1, &thing{0, 0, OtherX}))
	bad(ToExec(0, 0, &thing{0, 0, NoMode}))
	bad(ToExec(1, 0, &thing{0, 0, UserRWX}))
}
