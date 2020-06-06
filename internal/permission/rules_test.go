package permission

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_rules(t *testing.T) {
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	ok(ToRead(&Set{0, 0, UserR}, NewActor(0, 0)))
	ok(ToRead(&Set{0, 0, GroupR}, NewActor(0, 0)))
	ok(ToRead(&Set{0, 0, OtherR}, NewActor(1, 1)))
	bad(ToRead(&Set{0, 0, NoMode}, NewActor(0, 0)))
	bad(ToRead(&Set{0, 0, UserRWX}, NewActor(1, 0)))

	ok(ToWrite(&Set{0, 0, UserW}, NewActor(0, 0)))
	ok(ToWrite(&Set{0, 0, GroupW}, NewActor(0, 0)))
	ok(ToWrite(&Set{0, 0, OtherW}, NewActor(1, 1)))
	bad(ToWrite(&Set{0, 0, NoMode}, NewActor(0, 0)))
	bad(ToWrite(&Set{0, 0, UserRWX}, NewActor(1, 0)))

	ok(ToExec(&Set{0, 0, UserX}, NewActor(0, 0)))
	ok(ToExec(&Set{0, 0, GroupX}, NewActor(0, 0)))
	ok(ToExec(&Set{0, 0, OtherX}, NewActor(1, 1)))
	bad(ToExec(&Set{0, 0, NoMode}, NewActor(0, 0)))
	bad(ToExec(&Set{0, 0, UserRWX}, NewActor(1, 0)))

	ok(ToCreate(
		&Set{0, 0, UserRWX}, // parent
		&Set{0, 0, UserRWX}, // entity
		NewActor(0, 0),
	))
	ok(ToCreate(
		&Set{1, 1, OtherW},
		&Set{0, 0, UserRWX},
		NewActor(0, 0),
	))
	bad(ToCreate(
		&Set{1, 1, UserRWX},
		&Set{0, 0, UserRWX},
		NewActor(0, 0),
	))

	ok(ToDelete(
		&Set{0, 0, UserW},
		&Set{0, 0, UserW},
		NewActor(0, 0),
	))
	bad(ToDelete(
		&Set{0, 0, UserW},
		&Set{0, 0, UserR},
		NewActor(0, 0),
	))

	ok(ToUpdate(
		&Set{0, 0, UserX},
		&Set{0, 0, UserW},
		NewActor(0, 0),
	))
	bad(ToUpdate(
		&Set{0, 0, UserR},
		&Set{0, 0, UserW},
		NewActor(0, 0),
	))
}
