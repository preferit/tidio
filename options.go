package tidio

import "github.com/gregoryv/fox"

type ServiceOption interface {
	ForService(srv *Service) error
}

type UseLogger struct {
	fox.Logger
}

func (me UseLogger) ForService(srv *Service) error {
	srv.SetLogger(me)
	return nil
}

type InitialAccount struct {
	name   string
	secret string
}

func (me InitialAccount) ForService(srv *Service) error {
	return srv.AddAccount(me.name, me.secret)
}
