package tidio

import (
	"io"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_timesheets(t *testing.T) {
	var (
		assert  = asserter.New(t)
		ok, bad = assert().Errors()

		sheets = MemSheets{}.New()
		empty  = "{}"
	)
	sheets.Source = StringSource(empty)
	ok(sheets.Load())
	ok(sheets.Map(func(next *bool, s *Timesheet) error { return nil }))
	ok(sheets.AddTimesheet(Timesheet{}.New()))

	sheets.Destination = BrokenDestination{}
	bad(sheets.Save())
	bad(sheets.Map(func(next *bool, s *Timesheet) error { return io.EOF }))
}
