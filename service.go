package tidio

import (
	"fmt"

	"github.com/gregoryv/box"
)

// NewService returns a service with preconfigured options.
// Options may be *Store or APIKeys
func NewService(options ...interface{}) *Service {
	service := &Service{
		state: NewState(),
	}
	for _, opt := range options {
		switch opt := opt.(type) {
		case Accounts:
			service.Accounts = opt
		default:
			panic(fmt.Sprintf("%T not recognized", opt))
		}
	}
	return service
}

type Service struct {
	datafile string // where data is saved
	state    *State

	Accounts
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
