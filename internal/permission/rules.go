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
	o := e.PermSet()
	switch {
	case owner(e, a) && (o.Mode&UserR == UserR):
	case member(e, a) && (o.Mode&GroupR == GroupR):
	case o.Mode&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToWrite(e Secured, a *Actor) error {
	o := e.PermSet()
	switch {
	case owner(e, a) && (o.Mode&UserW == UserW):
	case member(e, a) && (o.Mode&GroupW == GroupW):
	case o.Mode&OtherW == OtherW:
	default:
		return ErrDenied
	}
	return nil
}

func (Rules) ToExec(e Secured, a *Actor) error {
	o := e.PermSet()
	switch {
	case owner(e, a) && (o.Mode&UserX == UserX):
	case member(e, a) && (o.Mode&GroupX == GroupX):
	case o.Mode&OtherX == OtherX:
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
	return a.UID == e.PermSet().UID
}

func member(e Secured, a *Actor) bool {
	for _, GID := range a.Groups {
		if GID == e.PermSet().GID {
			return true
		}
	}
	return false
}

var ErrDenied = errors.New("permission denied")
