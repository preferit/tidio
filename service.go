package tidio

import (
	"fmt"

	"github.com/gregoryv/box"
)

// NewService returns a service with preconfigured options.
// Options may be *Store or APIKeys
func NewService(options ...interface{}) *Service {
	service := &Service{
		data: &Data{},
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
	data     *Data

	apikeys APIKeys
}

func (s *Service) LoadState(filename string) error {
	s.datafile = filename
	return s.data.Load(&box.Store{}, filename)
}

func (s *Service) SaveState() error {
	return s.data.Save(&box.Store{}, s.datafile)
}

type APIKeys map[string]*Account

func (s *Service) IsAuthenticated(key string) (*Role, bool) {
	if key == "" {
		return nil, false
	}
	account, found := s.apikeys[key]
	if !found {
		return nil, false
	}
	return &Role{
		account: account,
		data:    s.data,
	}, true
}
