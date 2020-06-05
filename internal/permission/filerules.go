package permission

type FileRules struct{}

func (FileRules) ToRead(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserR == UserR):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupR == GroupR):
	case e.Mode()&OtherR == OtherR:
	default:
		return ErrDenied
	}
	return nil
}

func (FileRules) ToWrite(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserW == UserW):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupW == GroupW):
	case e.Mode()&OtherW == OtherW:
	default:
		return ErrDenied
	}
	return nil
}

func (FileRules) ToExec(e Resource, a Account) error {
	switch {
	case e.UID() == a.UID() && (e.Mode()&UserX == UserX):
	case a.Member(e.GID()) == nil && (e.Mode()&GroupX == GroupX):
	case e.Mode()&OtherX == OtherX:
	default:
		return ErrDenied
	}
	return nil
}
