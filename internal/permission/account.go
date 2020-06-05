package permission

import "errors"

type Account interface {
	UID() int
	Member(gid int) error
}

var ErrMembership = errors.New("missing membership")
