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

func ToWrite(uid, gid int, e Resource) error {
	switch {
	case e.UID() == uid && (e.Mode()&UserW == UserW):
	case e.GID() == gid && (e.Mode()&GroupW == GroupW):
	case e.Mode()&OtherW == OtherW:
	default:
		return ErrDenied
	}
	return nil
}

func ToExec(uid, gid int, e Resource) error {
	switch {
	case e.UID() == uid && (e.Mode()&UserX == UserX):
	case e.GID() == gid && (e.Mode()&GroupX == GroupX):
	case e.Mode()&OtherX == OtherX:
	default:
		return ErrDenied
	}
	return nil
}

var ErrDenied = errors.New("permission denied")
