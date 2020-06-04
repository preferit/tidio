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
	asAdmin := &Role{
		account: NewAccount("john", "admin"),
		store:   store,
	}

	t.Run("CreateTimesheet", func(t *testing.T) {
		assert := asserter.New(t)
		ok, bad := assert().Errors()

		bad(asAdmin.CreateTimesheet(&Timesheet{
			Filename: "xx.txt",
			Owner:    "john",
			Content:  aFile("x"),
		}))

		ok(asAdmin.CreateTimesheet(&Timesheet{
			Filename: "202001.timesheet",
			Owner:    "john",
			Content:  aFile("."),
		}))

		bad(asAdmin.CreateTimesheet(&Timesheet{
			Filename: "202001.timesheet",
			Owner:    "not-user",
			Content:  aFile("."),
		}))
	})

	t.Run("ReadTimesheet", func(t *testing.T) {
		filename := "199902.timesheet"
		s := &Timesheet{
			Filename: filename,
			Owner:    "john",
			Content:  aFile("..."),
		}
		asAdmin.CreateTimesheet(s)
		err := asAdmin.ReadTimesheet(ioutil.Discard, filename, "john")
		if err != nil {
			t.Error(err)
		}
		err = asAdmin.ReadTimesheet(ioutil.Discard, filename, "unknown")
		if err == nil {
			t.Error("read non existing file")
		}
	})

	t.Run("ListTimesheet", func(t *testing.T) {
		// depends on above tests creating some
		sheets := asAdmin.ListTimesheet("john")
		if len(sheets) == 0 {
			t.Error("did not found any timesheets")
		}
	})
}

func aFile(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
