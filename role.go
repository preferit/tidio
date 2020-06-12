package tidio

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Role struct {
	account *Account
	sheets  *MemSheets
}

func (r *Role) CreateTimesheet(sheet *Timesheet) error {
	if err := checkTimesheetFilename(sheet.Filename); err != nil {
		return err
	}
	var sb strings.Builder
	io.Copy(&sb, sheet)
	sheet.Content = sb.String()
	return r.sheets.AddTimesheet(sheet)
}

func (r *Role) OpenTimesheet(sheet *Timesheet) error {
	for _, s := range r.sheets.Sheets {
		if s.Equal(sheet) {
			*sheet = *s
			sheet.ReadCloser = ioutil.NopCloser(strings.NewReader(sheet.Content))
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func (r *Role) ListTimesheet(user string) []string {
	res := make([]string, 0)
	for _, s := range r.sheets.Sheets {
		if s.Owner == user {
			res = append(res, s.Filename)
		}
	}
	return res
}
