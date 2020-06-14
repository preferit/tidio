package tidio

import (
	"path"
)

func NewService() *Service {
	return &Service{
		Timesheets: NewMemSheets(),
		Accounts:   NewMemAccounts(),
	}
}

type Service struct {
	Stateful
	Timesheets
	Accounts
}

func (s *Service) Load() error {
	err := errors{
		s.Timesheets.Load(),
		s.Accounts.Load(),
	}
	return err.First()
}

func (s *Service) Save() error {
	err := errors{
		s.Timesheets.Save(),
		s.Accounts.Save(),
	}
	return err.First()
}

// SetDataDir sets directory where state is persisted
func (s *Service) SetDataDir(dir string) {
	s.Timesheets.PersistToFile(path.Join(dir, "timesheets.json"))
	s.Accounts.PersistToFile(path.Join(dir, "accounts.json"))
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
