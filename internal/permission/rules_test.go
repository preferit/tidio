package permission

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_rules(t *testing.T) {
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	ok(ToRead(&SecInfo{0, 0, UserR}, &mockAccount{0, []int{0}}))
	ok(ToRead(&SecInfo{0, 0, GroupR}, &mockAccount{0, []int{0}}))
	ok(ToRead(&SecInfo{0, 0, OtherR}, &mockAccount{1, []int{1}}))
	bad(ToRead(&SecInfo{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToRead(&SecInfo{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToWrite(&SecInfo{0, 0, UserW}, &mockAccount{0, []int{0}}))
	ok(ToWrite(&SecInfo{0, 0, GroupW}, &mockAccount{0, []int{0}}))
	ok(ToWrite(&SecInfo{0, 0, OtherW}, &mockAccount{1, []int{1}}))
	bad(ToWrite(&SecInfo{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToWrite(&SecInfo{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToExec(&SecInfo{0, 0, UserX}, &mockAccount{0, []int{0}}))
	ok(ToExec(&SecInfo{0, 0, GroupX}, &mockAccount{0, []int{0}}))
	ok(ToExec(&SecInfo{0, 0, OtherX}, &mockAccount{1, []int{1}}))
	bad(ToExec(&SecInfo{0, 0, NoMode}, &mockAccount{0, []int{0}}))
	bad(ToExec(&SecInfo{0, 0, UserRWX}, &mockAccount{1, []int{0}}))

	ok(ToCreate(
		&SecInfo{0, 0, UserRWX}, // parent
		&SecInfo{0, 0, UserRWX}, // entity
		newAccount(0, 0),
	))
	ok(ToCreate(
		&SecInfo{1, 1, OtherW},
		&SecInfo{0, 0, UserRWX},
		newAccount(0, 0),
	))
	bad(ToCreate(
		&SecInfo{1, 1, UserRWX},
		&SecInfo{0, 0, UserRWX},
		newAccount(0, 0),
	))

	ok(ToDelete(
		&SecInfo{0, 0, UserW},
		&SecInfo{0, 0, UserW},
		newAccount(0, 0),
	))
	bad(ToDelete(
		&SecInfo{0, 0, UserW},
		&SecInfo{0, 0, UserR},
		newAccount(0, 0),
	))

	ok(ToUpdate(
		&SecInfo{0, 0, UserX},
		&SecInfo{0, 0, UserW},
		newAccount(0, 0),
	))
	bad(ToUpdate(
		&SecInfo{0, 0, UserR},
		&SecInfo{0, 0, UserW},
		newAccount(0, 0),
	))
}
