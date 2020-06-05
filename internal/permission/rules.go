package permission

import "errors"

func ToRead(uid, gid int, e Resource) error {
	switch {
	case e.UID() == uid && (e.Mode()&UserR == UserR):
	case e.GID() == gid && (e.Mode()&GroupR == GroupR):
	case e.Mode()&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

var ErrDenied = errors.New("permission denied")
