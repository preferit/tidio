package tidio

import (
	"io"
	"strings"
)

type Role struct {
	account *Account
	Timesheets
}

func (r *Role) CreateTimesheet(sheet *Timesheet) error {
	if err := checkTimesheetFilename(sheet.FileSource); err != nil {
		return err
	}
	var sb strings.Builder
	io.Copy(&sb, sheet)
	sheet.Content = sb.String()
	return r.AddTimesheet(sheet)
}

func (r *Role) OpenTimesheet(sheet *Timesheet) error {
	return r.FindTimesheet(sheet)
}

func (r *Role) ListTimesheet() []string {
	res := make([]string, 0)
	r.Timesheets.Map(func(next *bool, s *Timesheet) error {
		// todo use account as filter
		res = append(res, s.FileSource)
		return nil
	})
	return res
}
