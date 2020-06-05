package permission

import "errors"

func ToRead(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserR == UserR):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupR == GroupR):
	case e.Mode()&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

func ToWrite(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserW == UserW):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupW == GroupW):
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
