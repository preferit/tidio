package tidio

import "github.com/gregoryv/fox"

type ServiceOptFunc func(srv *Service) error

func (me ServiceOptFunc) ForService(srv *Service) error {
	return me(srv)
}

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

func (me UseLogger) ForClient(cli *Client) error {
	cli.Logger = me
	return nil
}

type InitialAccount struct {
	account string
	secret  string
}

func (me InitialAccount) ForService(srv *Service) error {
	return srv.AddAccount(me.account, me.secret)
}

type UseHost string

func (me UseHost) ForClient(cli *Client) error {
	cli.host = string(me)
	return nil
}
