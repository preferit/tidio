package tidio

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/gregoryv/nugo"
	"github.com/preferit/tidio/internal"
)

func NewAccount(username string) *Account {
	return &Account{
		Username:   username,
		Timesheets: NewMemSheets(),
	}
}

type Account struct {
	nugo.Ring
	Username string
	Key      string

	Timesheets `json:"-"`
}

func (me *Account) WriteResource(resource *Resource) error {
	if err := isTimesheet(resource.Path); err != nil {
		return err
	}
	if err := checkTimesheetFilename(resource.Path); err != nil {
		return err
	}
	var sb strings.Builder
	io.Copy(&sb, resource)
	sheet := Timesheet{
		Path:    resource.Path,
		Content: sb.String(),
	}
	return me.AddTimesheet(&sheet)
}

func (me *Account) ReadResource(resource *Resource) error {
	sheet := Timesheet{
		Path: resource.Path,
	}
	if err := me.FindTimesheet(&sheet); err != nil {
		return err
	}
	resource.ReadCloser = sheet.ReadCloser
	return nil
}

func isTimesheet(filename string) error {
	if path.Ext(filename) != ".timesheet" {
		return fmt.Errorf("isTimesheet: %s", filename)
	}
	return nil
}

func (me *Account) OpenTimesheet(sheet *Timesheet) error {
	return me.FindTimesheet(sheet)
}

func (me *Account) ListTimesheet() []string {
	res := make([]string, 0)
	me.Timesheets.Map(func(next *bool, s *Timesheet) error {
		// todo use account as filter
		res = append(res, s.Path)
		return nil
	})
	return res
}

// ----------------------------------------

type Accounts interface {
	Stateful
	FilePersistent
	AddAccount(*Account) error
	FindAccountByKey(*Account, string) error
}

// ----------------------------------------

func NewMemAccounts() *MemAccounts {
	return &MemAccounts{
		accounts:    make(map[string]*Account),
		Source:      internal.None("AccountsMap"),
		Destination: internal.None("AccountsMap"),
	}
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

func (s *MemAccounts) AddAccount(a *Account) error {
	// todo key should not be required
	s.accounts[a.Key] = a
	return nil
}

func (s *MemAccounts) FindAccountByKey(a *Account, key string) error {
	for _, account := range s.accounts {
		if account.Key == key {
			*a = *account
			return nil
		}
	}
	return fmt.Errorf("account not found")
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
