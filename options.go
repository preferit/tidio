package tidio

import "github.com/gregoryv/fox"

type ServiceOption interface {
	ForService(srv *Service) error
}

type LoggerOption struct {
	fox.Logger
}

func (me LoggerOption) ForService(srv *Service) error {
	srv.SetLogger(me)
	return nil
}
