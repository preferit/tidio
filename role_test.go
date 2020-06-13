package tidio

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_role(t *testing.T) {
	var (
		assert  = asserter.New(t)
		ok, bad = assert().Errors()
		john    = &Role{
			account:    NewAccount("john", "admin"),
			Timesheets: &MemSheets{},
		}
	)

	bad(john.CreateTimesheet(&Timesheet{
		FileSource: "xx.txt",
		Owner:      "john",
		ReadCloser: aFile("x")},
	))

	ok(john.CreateTimesheet(&Timesheet{
		FileSource: "202001.timesheet",
		Owner:      "john",
		ReadCloser: aFile("."),
	}))
	ok(john.OpenTimesheet(&Timesheet{
		FileSource: "202001.timesheet",
		Owner:      "john",
	}))

	bad(john.OpenTimesheet(&Timesheet{
		FileSource: "209901.timesheet",
		Owner:      "john",
	}))
	t.Run("ListTimesheet", func(t *testing.T) {
		// depends on above tests creating some
		Sheets := john.ListTimesheet("john")
		if len(Sheets) == 0 {
			t.Error("did not found any timesheets")
		}
	})
}

func aFile(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
