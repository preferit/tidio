package permission

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_account(t *testing.T) {
	assert := asserter.New(t)
	ok, bad := assert().Errors()

	var account Account
	account = &mockAccount{uid: 0, groups: []int{0, 2}}
	account.UID()

	ok(account.Member(0))
	bad(account.Member(9))
}

func newAccount(uid int, groups ...int) *mockAccount {
	return &mockAccount{
		uid:    uid,
		groups: groups,
	}
}

type mockAccount struct {
	uid    int
	groups []int
}

func (a *mockAccount) UID() int { return a.uid }
func (a *mockAccount) Member(gid int) error {
	for _, group := range a.groups {
		if gid == group {
			return nil
		}
	}
	return ErrMembership
}
