package tidio

import (
	"io"
	"os"
)

type Service struct {
	datafile string // where data is saved
	Timesheets
	Accounts
}

func (s Service) New() *Service {
	e := &s
	e.Timesheets = MemSheets{}.New()
	return e
}

func (s *Service) LoadState(filename string) error {
	s.datafile = filename
	return s.Timesheets.ReadState(func() (io.ReadCloser, error) {
		return os.Open(filename)
	})
}

func (s *Service) SaveState() error {
	return s.Timesheets.WriteState(os.Create(s.datafile))
}

func (s *Service) RoleByKey(key string) (*Role, bool) {
	if key == "" {
		return nil, false
	}
	var account Account
	if err := s.FindAccountByKey(&account, key); err != nil {
		return nil, false
	}
	return &Role{
		account:    &account,
		Timesheets: s.Timesheets,
	}, true
}
