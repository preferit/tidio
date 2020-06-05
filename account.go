package tidio

import "github.com/preferit/tidio/internal/permission"

func NewAccount(username, role string) *Account {
	uid := 0
	gid := 0
	return &Account{
		UID_:     uid,
		Groups:   []int{gid},
		Username: username,
		Role:     role,
	}
}

type Account struct {
	UID_     int // _ suffix so we can implement permission.Account
	Groups   []int
	Username string
	Role     string
}

func (a *Account) UID() int { return a.UID_ }
func (a *Account) Member(gid int) error {
	for _, group := range a.Groups {
		if group == gid {
			return nil
		}
	}
	return permission.ErrMembership
}
