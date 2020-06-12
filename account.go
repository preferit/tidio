package tidio

import "fmt"

func NewAccount(username, role string) *Account {
	return &Account{
		Username: username,
	}
}

type Account struct {
	Username string
}

// ----------------------------------------

type Accounts interface {
	FindAccountByKey(*Account, string) error
}

// ----------------------------------------

type AccountsMap map[string]*Account

func (s AccountsMap) FindAccountByKey(a *Account, key string) error {
	account, found := s[key]
	if !found {
		return fmt.Errorf("account not found")
	}
	*a = *account
	return nil
}
