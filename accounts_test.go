package tidio

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_accounts(t *testing.T) {
	var (
		assert   = asserter.New(t)
		ok, bad  = assert().Errors()
		empty    = "{}"
		accounts = NewMemAccounts()
		acc      Account
	)
	accounts.Source = StringSource(empty)
	accounts.Destination = Nowhere()

	ok(accounts.Load())
	ok(accounts.Save())

	accounts.Source = BrokenSource{}
	accounts.Destination = BrokenDestination{}

	bad(accounts.FindAccountByKey(&acc, "x"))
	bad(accounts.Load())
	bad(accounts.Save())
}

func Test_account(t *testing.T) {
	var (
		assert  = asserter.New(t)
		ok, bad = assert().Errors()
		john    = NewAccount("john")
	)
	john.Timesheets = NewMemSheets()

	ok(john.WriteResource(&Resource{
		Path:       "202001.timesheet",
		ReadCloser: aFile("."),
	}))
	ok(john.OpenTimesheet(&Timesheet{
		Path: "202001.timesheet",
	}))

	bad(john.OpenTimesheet(&Timesheet{
		Path: "209901.timesheet",
	}))
	t.Run("FindResources", func(t *testing.T) {
		Sheets := john.FindResources()
		if len(Sheets) == 0 {
			t.Error("did not found any timesheets")
		}
	})
}

func aFile(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
