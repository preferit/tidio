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

func (c *Rules) ToUpdate(parent, e Resource, a Account) error {
	if c.ToExec(parent, a) != nil || c.ToWrite(e, a) != nil {
		return ErrDenied
	}
	return nil
}

func (c *Rules) ToDelete(parent, e Resource, a Account) error {
	if c.ToWrite(parent, a) != nil || c.ToWrite(e, a) != nil {
		return ErrDenied
	}
	return nil
}

func (Rules) ToRead(e Resource, a Account) error {
	o := e.SecInfo()
	switch {
	case owner(e, a) && (o.mode&UserR == UserR):
	case a.Member(o.gid) == nil && (o.mode&GroupR == GroupR):
	case o.mode&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToWrite(e Resource, a Account) error {
	o := e.SecInfo()
	switch {
	case owner(e, a) && (o.mode&UserW == UserW):
	case a.Member(o.gid) == nil && (o.mode&GroupW == GroupW):
	case o.mode&OtherW == OtherW:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToExec(e Resource, a Account) error {
	o := e.SecInfo()
	switch {
	case owner(e, a) && (o.mode&UserX == UserX):
	case a.Member(o.gid) == nil && (o.mode&GroupX == GroupX):
	case o.mode&OtherX == OtherX:
	default:
		return ErrDenied
	}
	return nil
}

func ToCreate(parent, e Resource, a Account) error {
	return DefaultRules.ToCreate(parent, e, a)
}

func ToDelete(parent, e Resource, a Account) error {
	return DefaultRules.ToDelete(parent, e, a)
}

func ToUpdate(parent, e Resource, a Account) error {
	return DefaultRules.ToUpdate(parent, e, a)
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
	return a.UID() == e.SecInfo().uid
}

var ErrDenied = errors.New("permission denied")
