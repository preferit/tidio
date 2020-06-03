package tidio

import (
	"io"
	"path"
)

type Role struct {
	account string
	store   *Store
}

func (r *Role) Account() string {
	return r.account
}

func (r *Role) CreateTimesheet(filename, user string, content io.ReadCloser) error {
	if err := checkTimesheetFilename(filename); err != nil {
		return err
	}
	if user != r.Account() {
		return ErrForbidden
	}
	filename = path.Join(user, filename)
	return r.store.WriteFile(filename, content)
}

func (r *Role) ReadTimesheet(w io.Writer, filename, user string) error {
	// todo Role implementation should not have permissions checks
	// move to eg. admin
	if user != r.Account() {
		return ErrForbidden
	}
	filename = path.Join(user, filename)
	return r.store.ReadFile(w, filename)
}
