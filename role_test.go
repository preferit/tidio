package tidio

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_role(t *testing.T) {
	john := &Role{
		account: NewAccount("john", "admin"),
		sheets:  &MemSheets{},
	}

	assert := asserter.New(t)
	ok, bad := assert().Errors()

	bad(john.CreateTimesheet(&Timesheet{
		Filename:   "xx.txt",
		Owner:      "john",
		ReadCloser: aFile("x")},
	))

	ok(john.CreateTimesheet(&Timesheet{
		Filename:   "202001.timesheet",
		Owner:      "john",
		ReadCloser: aFile("."),
	}))

	ok(john.OpenTimesheet(&Timesheet{
		Filename: "202001.timesheet",
		Owner:    "john",
	}))

	bad(john.OpenTimesheet(&Timesheet{
		Filename: "209901.timesheet",
		Owner:    "john",
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
