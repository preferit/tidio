package tidio

import "github.com/gregoryv/fox"

type ServiceOption interface {
	ForService(srv *Service) error
}

type ClientOption interface {
	ForClient(cli *Client) error
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

type UseHost string

func (me UseHost) ForClient(cli *Client) error {
	cli.host = string(me)
	return nil
}
