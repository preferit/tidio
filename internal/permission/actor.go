package permission

import "errors"

func NewActor(uid int, groups ...int) *Actor {
	return &Actor{
		UID:    uid,
		Groups: groups,
	}
}

type Actor struct {
	UID    int
	Groups []int
}

var ErrMembership = errors.New("missing membership")
