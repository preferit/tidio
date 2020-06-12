package tidio

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_timesheets(t *testing.T) {
	var (
		assert  = asserter.New(t)
		ok, bad = assert().Errors()

		sheets     Timesheets = MemSheets{}.New()
		WriteState            = sheets.WriteState
		ReadState             = sheets.ReadState
		Map                   = sheets.Map
		empty                 = "{}"
	)
	ok(WriteState(&nopWriteCloser{}, nil))
	bad(WriteState(nil, io.EOF))
	ok(ReadState(nopRead(empty), nil))
	ok(Map(func(next *bool, s *Timesheet) error { return nil }))

	ok(sheets.AddTimesheet(Timesheet{}.New()))
	bad(Map(func(next *bool, s *Timesheet) error { return io.EOF }))
}

func nopRead(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}
