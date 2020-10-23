package tidio

import (
	"github.com/gregoryv/ant"
	"github.com/gregoryv/fox"
)

type ErrorHandling func(...interface{})

func (me ErrorHandling) Set(v interface{}) error {
	switch v := v.(type) {
	case *Client:
		v.check = me
	default:
		return ant.SetFailed(v, me)
	}
	return nil
}

type Logging struct {
	fox.Logger
}

func (me Logging) Set(v interface{}) error {
	switch v := v.(type) {
	case usesLogging:
		v.SetLogger(me.Logger)
	default:
		return ant.SetFailed(v, me)
	}
	return nil
}

type usesLogging interface {
	SetLogger(fox.Logger)
}

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
