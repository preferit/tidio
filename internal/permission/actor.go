package permission

import "errors"

func NewActor(uid uint, groups ...uint) *Actor {
	return &Actor{
		UID:    uid,
		Groups: groups,
	}
}

type Actor struct {
	UID    uint
	Groups []uint
}

var ErrMembership = errors.New("missing membership")
