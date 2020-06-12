package tidio

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_accounts(t *testing.T) {
	var (
		assert                    = asserter.New(t)
		ok, bad                   = assert().Errors()
		accounts         Accounts = AccountsMap{}.New()
		FindAccountByKey          = accounts.FindAccountByKey
		LoadAccounts              = accounts.LoadAccounts
		SaveAccounts              = accounts.SaveAccounts
		acc              Account
		empty            = "{}"
	)
	bad(FindAccountByKey(&acc, "x"))
	ok(LoadAccounts(ioutil.NopCloser(strings.NewReader(empty))))
	ok(SaveAccounts(&nopWriteCloser{}))
}

type nopWriteCloser strings.Builder

func (nopWriteCloser) Write(b []byte) (int, error) {
	return ioutil.Discard.Write(b)
}

func (nopWriteCloser) Close() error {
	return nil
}
