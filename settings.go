package tidio

import (
	"github.com/gregoryv/ant"
)

// Not really a setting but very helpful
type InitialAccount struct {
	account string
	secret  string
}

func (me InitialAccount) Set(v interface{}) error {
	switch v := v.(type) {
	case *Service:
		v.AddAccount(me.account, me.secret)
	default:
		return ant.SetFailed(v, me)
	}
	return nil
}
