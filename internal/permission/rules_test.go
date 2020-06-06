package permission

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_rules(t *testing.T) {
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	ok(ToRead(&SecInfo{0, 0, UserR}, NewActor(0, 0)))
	ok(ToRead(&SecInfo{0, 0, GroupR}, NewActor(0, 0)))
	ok(ToRead(&SecInfo{0, 0, OtherR}, NewActor(1, 1)))
	bad(ToRead(&SecInfo{0, 0, NoMode}, NewActor(0, 0)))
	bad(ToRead(&SecInfo{0, 0, UserRWX}, NewActor(1, 0)))

	ok(ToWrite(&SecInfo{0, 0, UserW}, NewActor(0, 0)))
	ok(ToWrite(&SecInfo{0, 0, GroupW}, NewActor(0, 0)))
	ok(ToWrite(&SecInfo{0, 0, OtherW}, NewActor(1, 1)))
	bad(ToWrite(&SecInfo{0, 0, NoMode}, NewActor(0, 0)))
	bad(ToWrite(&SecInfo{0, 0, UserRWX}, NewActor(1, 0)))

	ok(ToExec(&SecInfo{0, 0, UserX}, NewActor(0, 0)))
	ok(ToExec(&SecInfo{0, 0, GroupX}, NewActor(0, 0)))
	ok(ToExec(&SecInfo{0, 0, OtherX}, NewActor(1, 1)))
	bad(ToExec(&SecInfo{0, 0, NoMode}, NewActor(0, 0)))
	bad(ToExec(&SecInfo{0, 0, UserRWX}, NewActor(1, 0)))

	ok(ToCreate(
		&SecInfo{0, 0, UserRWX}, // parent
		&SecInfo{0, 0, UserRWX}, // entity
		NewActor(0, 0),
	))
	ok(ToCreate(
		&SecInfo{1, 1, OtherW},
		&SecInfo{0, 0, UserRWX},
		NewActor(0, 0),
	))
	bad(ToCreate(
		&SecInfo{1, 1, UserRWX},
		&SecInfo{0, 0, UserRWX},
		NewActor(0, 0),
	))

	ok(ToDelete(
		&SecInfo{0, 0, UserW},
		&SecInfo{0, 0, UserW},
		NewActor(0, 0),
	))
	bad(ToDelete(
		&SecInfo{0, 0, UserW},
		&SecInfo{0, 0, UserR},
		NewActor(0, 0),
	))

	ok(ToUpdate(
		&SecInfo{0, 0, UserX},
		&SecInfo{0, 0, UserW},
		NewActor(0, 0),
	))
	bad(ToUpdate(
		&SecInfo{0, 0, UserR},
		&SecInfo{0, 0, UserW},
		NewActor(0, 0),
	))
}
