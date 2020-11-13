package tidio

import (
	"github.com/gregoryv/ant"
)

func NewInitialAccount(cred *Credentials) *InitialAccount {
	return &InitialAccount{
		Account: cred.account,
		Secret:  cred.secret,
	}
}

type InitialAccount struct {
	Account string
	Secret  string
}

func (me InitialAccount) Set(v interface{}) error {
	switch v := v.(type) {
	case *Service:
		v.AddAccount(me.Account, me.Secret)
	default:
		return ant.SetFailed(v, me)
	}
	return nil
}
