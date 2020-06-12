package tidio

import (
	"encoding/json"
	"fmt"
	"io"
)

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
	Stateful
	AddAccount(string, *Account) error
	FindAccountByKey(*Account, string) error
}

// ----------------------------------------

type AccountsMap struct {
	accounts map[string]*Account
}

func (s AccountsMap) New() *AccountsMap {
	accounts := &s
	accounts.accounts = make(map[string]*Account)
	return accounts
}

func (s *AccountsMap) AddAccount(key string, a *Account) error {
	s.accounts[key] = a
	return nil
}

func (s *AccountsMap) FindAccountByKey(a *Account, key string) error {
	account, found := s.accounts[key]
	if !found {
		return fmt.Errorf("account not found")
	}
	*a = *account
	return nil
}

func (s *AccountsMap) ReadState(open ReadOpener) error {
	r, err := open()
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(&s.accounts)
}

func (s *AccountsMap) WriteState(fh io.WriteCloser, err error) error {
	if err != nil {
		return err
	}
	defer fh.Close()
	return json.NewEncoder(fh).Encode(s.accounts)
}
