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
		Path:       "xx.txt",
		Owner:      "john",
		ReadCloser: aFile("x")},
	))

	ok(john.CreateTimesheet(&Timesheet{
		Path:       "202001.timesheet",
		Owner:      "john",
		ReadCloser: aFile("."),
	}))
	ok(john.OpenTimesheet(&Timesheet{
		Path:  "202001.timesheet",
		Owner: "john",
	}))

	bad(john.OpenTimesheet(&Timesheet{
		Path:  "209901.timesheet",
		Owner: "john",
	}))
	t.Run("ListTimesheet", func(t *testing.T) {

		Sheets := john.ListTimesheet()
		if len(Sheets) == 0 {
			t.Error("did not found any timesheets")
		}
	})
}

func aFile(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
