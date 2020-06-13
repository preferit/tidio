package tidio

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_accounts(t *testing.T) {
	var (
		assert   = asserter.New(t)
		ok, bad  = assert().Errors()
		empty    = "{}"
		accounts = AccountsMap{}.New()
		acc      Account
	)
	accounts.Source = StringSource(empty)
	accounts.Destination = NopDestination()

	ok(accounts.Load())
	ok(accounts.Save())

	accounts.Source = BrokenSource{}
	accounts.Destination = BrokenDestination{}

	bad(accounts.FindAccountByKey(&acc, "x"))
	bad(accounts.Load())
	bad(accounts.Save())
}
