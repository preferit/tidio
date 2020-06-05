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

	ok(ToExec(&thing{0, 0, UserX}, &mockAccount{0, []int{0}}))
	ok(ToExec(&thing{0, 0, GroupX}, &mockAccount{0, []int{0}}))
	ok(ToExec(&thing{0, 0, OtherX}, &mockAccount{1, []int{1}}))
	bad(ToExec(&thing{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToExec(&thing{0, 0, UserRWX}, &mockAccount{1, []int{0}}))
}
