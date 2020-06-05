package permission

import "errors"

var DefaultRules = &Rules{}

type Rules struct{}

func (c *Rules) ToCreate(parent, e Resource, a Account) error {
	if c.ToWrite(parent, a) != nil || !owner(e, a) {
		return ErrDenied
	}
	return nil
}

func (Rules) ToRead(e Resource, a Account) error {
	switch {
	case owner(e, a) && (e.Mode()&UserR == UserR):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupR == GroupR):
	case e.Mode()&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToWrite(e Resource, a Account) error {
	switch {
	case owner(e, a) && (e.Mode()&UserW == UserW):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupW == GroupW):
	case e.Mode()&OtherW == OtherW:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToExec(e Resource, a Account) error {
	switch {
	case owner(e, a) && (e.Mode()&UserX == UserX):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupX == GroupX):
	case e.Mode()&OtherX == OtherX:
	default:
		return ErrDenied
	}
	return nil
}

func ToCreate(parent, e Resource, a Account) error {
	return DefaultRules.ToCreate(parent, e, a)
}

func ToRead(e Resource, a Account) error {
	return DefaultRules.ToRead(e, a)
}

func ToWrite(e Resource, a Account) error {
	return DefaultRules.ToWrite(e, a)
}

func ToExec(e Resource, a Account) error {
	return DefaultRules.ToExec(e, a)
}

func owner(e Resource, a Account) bool {
	return a.UID() == e.UID()
}

var ErrDenied = errors.New("permission denied")
