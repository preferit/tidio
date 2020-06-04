package tidio

import (
	"io"
	"path"
)

type Role struct {
	account *Account
	store   *Store
}

type Timesheet struct {
	Filename string
	Owner    string
	Content  io.ReadCloser
}

func (r *Role) CreateTimesheet(s *Timesheet) error {
	if err := checkTimesheetFilename(s.Filename); err != nil {
		return err
	}
	if s.Owner != r.account.Username {
		return ErrForbidden
	}
	out := path.Join(s.Owner, s.Filename)
	return r.store.WriteFile(r.account.Username, out, s.Content)
}

func (r *Role) ReadTimesheet(w io.Writer, filename, user string) error {
	// todo Role implementation should not have permissions checks
	// move to eg. admin
	if user != r.account.Username {
		return ErrForbidden
	}
	filename = path.Join(user, filename)
	return r.store.ReadFile(w, filename)
}

func (r *Role) ListTimesheet(user string) []string {
	return r.store.Glob(user, "*.timesheet")
}
