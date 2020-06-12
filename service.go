package tidio

import (
	"github.com/gregoryv/box"
)

type Service struct {
	datafile string // where data is saved
	state    *State

	Accounts
}

func (s Service) New() *Service {
	e := &s
	e.state = NewState()
	return e
}

func (s *Service) LoadState(filename string) error {
	s.datafile = filename
	return s.state.Load(&box.Store{}, filename)
}

func (s *Service) SaveState() error {
	return s.state.Save(&box.Store{}, s.datafile)
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
		account: &account,
		state:   s.state,
	}, true
}
