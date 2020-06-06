package permission

import "errors"

var DefaultRules = &Rules{}

type Rules struct{}

func (c *Rules) ToCreate(parent, e Secured, a *Actor) error {
	if c.ToWrite(parent, a) != nil || !owner(e, a) {
		return ErrDenied
	}
	return nil
}

func (c *Rules) ToUpdate(parent, e Secured, a *Actor) error {
	if c.ToExec(parent, a) != nil || c.ToWrite(e, a) != nil {
		return ErrDenied
	}
	return nil
}

func (c *Rules) ToDelete(parent, e Secured, a *Actor) error {
	if c.ToWrite(parent, a) != nil || c.ToWrite(e, a) != nil {
		return ErrDenied
	}
	return nil
}

func (Rules) ToRead(e Secured, a *Actor) error {
	o := e.SecInfo()
	switch {
	case owner(e, a) && (o.mode&UserR == UserR):
	case member(e, a) && (o.mode&GroupR == GroupR):
	case o.mode&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToWrite(e Secured, a *Actor) error {
	o := e.SecInfo()
	switch {
	case owner(e, a) && (o.mode&UserW == UserW):
	case member(e, a) && (o.mode&GroupW == GroupW):
	case o.mode&OtherW == OtherW:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToExec(e Secured, a *Actor) error {
	o := e.SecInfo()
	switch {
	case owner(e, a) && (o.mode&UserX == UserX):
	case member(e, a) && (o.mode&GroupX == GroupX):
	case o.mode&OtherX == OtherX:
	default:
		return ErrDenied
	}
	return nil
}

func ToCreate(parent, e Secured, a *Actor) error {
	return DefaultRules.ToCreate(parent, e, a)
}

func ToDelete(parent, e Secured, a *Actor) error {
	return DefaultRules.ToDelete(parent, e, a)
}

func ToUpdate(parent, e Secured, a *Actor) error {
	return DefaultRules.ToUpdate(parent, e, a)
}

func ToRead(e Secured, a *Actor) error {
	return DefaultRules.ToRead(e, a)
}

func ToWrite(e Secured, a *Actor) error {
	return DefaultRules.ToWrite(e, a)
}

func ToExec(e Secured, a *Actor) error {
	return DefaultRules.ToExec(e, a)
}

func owner(e Secured, a *Actor) bool {
	return a.UID == e.SecInfo().uid
}

func member(e Secured, a *Actor) bool {
	for _, gid := range a.Groups {
		if gid == e.SecInfo().gid {
			return true
		}
	}
	return false
}

var ErrDenied = errors.New("permission denied")
