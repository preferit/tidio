package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
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

var ErrForbidden = fmt.Errorf("forbidden")

func checkTimesheetFilename(name string) error {
	format := `\d\d\d\d\d\d\.timesheet`
	if ok, _ := regexp.MatchString(format, name); !ok {
		return fmt.Errorf("bad filename: expected format %s", format)
	}
	return nil
}
