package permission

import "errors"

var DefaultRules Rules = &FileRules{}

type Rules interface {
	ToRead(Resource, Account) error
	ToWrite(Resource, Account) error
	ToExec(Resource, Account) error
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

var ErrDenied = errors.New("permission denied")
