package permission

import "errors"

var DefaultRules = &Rules{}

func ToRead(e Resource, a Account) error {
	return DefaultRules.ToRead(e, a)
}

func ToWrite(e Resource, a Account) error {
	return DefaultRules.ToWrite(e, a)
}

func ToExec(e Resource, a Account) error {
	return DefaultRules.ToExec(e, a)
}

type Rules struct{}

func (*Rules) ToRead(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserR == UserR):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupR == GroupR):
	case e.Mode()&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

func (*Rules) ToWrite(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserW == UserW):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupW == GroupW):
	case e.Mode()&OtherW == OtherW:
	default:
		return ErrDenied
	}
	return nil
}

func (*Rules) ToExec(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserX == UserX):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupX == GroupX):
	case e.Mode()&OtherX == OtherX:
	default:
		return ErrDenied
	}
	return nil
}

var ErrDenied = errors.New("permission denied")
