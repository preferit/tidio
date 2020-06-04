package tidio

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_role(t *testing.T) {
	store, cleanup := newTempStore(t)
	defer cleanup()
	john := &Role{
		account: NewAccount("john", "admin"),
		store:   store,
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

	bad(john.CreateTimesheet(&Timesheet{
		Filename:   "202001.timesheet",
		Owner:      "not-user",
		ReadCloser: aFile("."),
	}))

	t.Run("ReadTimesheet", func(t *testing.T) {
		filename := "199902.timesheet"
		s := &Timesheet{
			Filename:   filename,
			Owner:      "john",
			ReadCloser: aFile("..."),
		}
		john.CreateTimesheet(s)
		err := john.OpenTimesheet(s)
		if err != nil {
			t.Error(err)
		}
		err = john.OpenTimesheet(&Timesheet{
			Filename: filename,
			Owner:    "unknown"})
		if err == nil {
			t.Error("read non existing file")
		}
	})

	t.Run("ListTimesheet", func(t *testing.T) {
		// depends on above tests creating some
		sheets := john.ListTimesheet("john")
		if len(sheets) == 0 {
			t.Error("did not found any timesheets")
		}
	})
}

func aFile(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
