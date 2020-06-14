package tidio

import (
	"encoding/json"
	"fmt"
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
	FilePersistent
	AddAccount(string, *Account) error
	FindAccountByKey(*Account, string) error
}

// ----------------------------------------

func NewMemAccounts() *MemAccounts {
	my := &MemAccounts{}
	my.accounts = make(map[string]*Account)
	my.Source = None("AccountsMap")
	my.Destination = None("AccountsMap")
	return my
}

type MemAccounts struct {
	Source
	Destination
	accounts map[string]*Account
}

func (me *MemAccounts) PersistToFile(filename string) {
	me.Source = FileSource(filename)
	me.Destination = FileDestination(filename)
}

func (s *MemAccounts) AddAccount(key string, a *Account) error {
	s.accounts[key] = a
	return nil
}

func (s *MemAccounts) FindAccountByKey(a *Account, key string) error {
	account, found := s.accounts[key]
	if !found {
		return fmt.Errorf("account not found")
	}
	*a = *account
	return nil
}

func (s *MemAccounts) Load() error {
	r, err := s.Source.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(&s.accounts)
}

func (s *MemAccounts) Save() error {
	w, err := s.Destination.Create()
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(s.accounts)
}
