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
		case APIKeys:
			service.apikeys = opt
		default:
			panic(fmt.Sprintf("%T not recognized", opt))
		}
	}
	return service
}

type Service struct {
	datafile string // where data is saved
	state    *State

	apikeys APIKeys
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
	account, found := s.apikeys[key]
	if !found {
		return nil, false
	}
	return &Role{
		account: account,
		state:   s.state,
	}, true
}
