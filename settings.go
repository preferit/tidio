package tidio

import (
	"fmt"

	"github.com/gregoryv/fox"
)

type SetFunc func(v interface{}) error

func (me SetFunc) Set(v interface{}) error {
	return me(v)
}

type Setting interface {
	Set(v interface{}) error
}

// ----------------------------------------

type ErrorHandling func(...interface{})

func (me ErrorHandling) Set(v interface{}) error {
	switch v := v.(type) {
	case *Client:
		v.check = me
	default:
		return setErr(me, v)
	}
	return nil
}

type UseLogger struct {
	fox.Logger
}

func (me UseLogger) Set(v interface{}) error {
	switch v := v.(type) {
	case *Service:
		v.SetLogger(me)
	case *Client:
		v.Logger = me
	default:
		return setErr(me, v)
	}
	return nil
}

type InitialAccount struct {
	account string
	secret  string
}

// Method
func (me InitialAccount) Set(v interface{}) error {
	switch v := v.(type) {
	case *Service:
		v.AddAccount(me.account, me.secret)
	default:
		return setErr(me, v)
	}
	return nil
}

type UseHost string

func (me UseHost) Set(v interface{}) error {
	switch v := v.(type) {
	case *Client:
		v.host = string(me)
	default:
		return setErr(me, v)
	}
	return nil
}

func setErr(s, v interface{}) error {
	return fmt.Errorf("%t cannot be set on %t", s, v)
}
