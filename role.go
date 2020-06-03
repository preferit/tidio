package tidio

import "io"

type Role struct {
	account string
	service *Service
}

func (r *Role) Account() string {
	return r.account
}

func (r *Role) CreateTimesheet(filename string, content io.ReadCloser) error {
	return nil
}
