package tidio

import "path"

type Service struct {
	Stateful

	Timesheets
	Accounts
}

func (s Service) New() *Service {
	e := &s
	e.Timesheets = MemSheets{}.New()
	e.Accounts = AccountsMap{}.New()
	return e
}

func (s *Service) SetDataDir(dir string) {
	s.Timesheets.PersistToFile(path.Join(dir, "timesheets.json"))
	s.Accounts.PersistToFile(path.Join(dir, "accounts.json"))
}

func (s *Service) Load() error {
	return s.Timesheets.Load()
}

func (s *Service) Save() error {
	return s.Timesheets.Save()
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
