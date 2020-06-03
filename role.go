package tidio

import (
	"io"
	"io/ioutil"
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
	body, _ := ioutil.ReadAll(content)
	return r.store.WriteFile(filename, body, 0644)
}

func (r *Role) ReadTimesheet(w io.Writer, filename, user string) error {
	return nil
}
