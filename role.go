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
	io.ReadCloser
	Content string
}

func (s *Timesheet) SetStream(stream io.ReadCloser) {
	s.ReadCloser = stream
}

func (r *Role) CreateTimesheet(sheet *Timesheet) error {
	if err := checkTimesheetFilename(sheet.Filename); err != nil {
		return err
	}
	r.store.Add(sheet)
	out := path.Join(sheet.Owner, sheet.Filename)
	return r.store.WriteFile(r.account.Username, out, sheet)
}

func (r *Role) OpenTimesheet(sheet *Timesheet) error {

	filename := path.Join(sheet.Owner, sheet.Filename)
	return r.store.OpenFile(sheet, filename)
}

func (r *Role) ListTimesheet(user string) []string {
	return r.store.Glob(user, "*.timesheet")
}
