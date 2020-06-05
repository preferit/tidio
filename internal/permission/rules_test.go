package permission

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_rules(t *testing.T) {
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	ok(ToRead(&Info{0, 0, UserR}, &mockAccount{0, []int{0}}))
	ok(ToRead(&Info{0, 0, GroupR}, &mockAccount{0, []int{0}}))
	ok(ToRead(&Info{0, 0, OtherR}, &mockAccount{1, []int{1}}))
	bad(ToRead(&Info{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToRead(&Info{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToWrite(&Info{0, 0, UserW}, &mockAccount{0, []int{0}}))
	ok(ToWrite(&Info{0, 0, GroupW}, &mockAccount{0, []int{0}}))
	ok(ToWrite(&Info{0, 0, OtherW}, &mockAccount{1, []int{1}}))
	bad(ToWrite(&Info{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToWrite(&Info{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToExec(&Info{0, 0, UserX}, &mockAccount{0, []int{0}}))
	ok(ToExec(&Info{0, 0, GroupX}, &mockAccount{0, []int{0}}))
	ok(ToExec(&Info{0, 0, OtherX}, &mockAccount{1, []int{1}}))
	bad(ToExec(&Info{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToExec(&Info{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToCreate(
		&Info{0, 0, UserRWX}, // parent
		&Info{0, 0, UserRWX}, // entity
		newAccount(0, 0),
	))
	ok(ToCreate(
		&Info{1, 1, OtherW},
		&Info{0, 0, UserRWX},
		newAccount(0, 0),
	))
	bad(ToCreate(
		&Info{1, 1, UserRWX},
		&Info{0, 0, UserRWX},
		newAccount(0, 0),
	))

	ok(ToDelete(
		&Info{0, 0, UserW},
		&Info{0, 0, UserW},
		newAccount(0, 0),
	))
	bad(ToDelete(
		&Info{0, 0, UserW},
		&Info{0, 0, UserR},
		newAccount(0, 0),
	))

	ok(ToUpdate(
		&Info{0, 0, UserX},
		&Info{0, 0, UserW},
		newAccount(0, 0),
	))
	bad(ToUpdate(
		&Info{0, 0, UserR},
		&Info{0, 0, UserW},
		newAccount(0, 0),
	))
}
