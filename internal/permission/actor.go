package permission

import "errors"

func NewActor(UID uint, groups ...uint) *Actor {
	return &Actor{
		UID:    UID,
		Groups: groups,
	}
}

type Actor struct {
	UID    uint
	Groups []uint
}

var ErrMembership = errors.New("missing membership")
