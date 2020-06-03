package tidio

import "fmt"

// NewService returns a service with preconfigured options.
// Options may be *Store or APIKeys
func NewService(options ...interface{}) *Service {
	service := &Service{}
	for _, opt := range options {
		switch opt := opt.(type) {
		case *Store:
			service.store = opt
		case APIKeys:
			service.apikeys = opt
		default:
			panic(fmt.Sprintf("%T not recognized", opt))
		}
	}
	return service
}

type Service struct {
	store   *Store
	apikeys APIKeys
}

type APIKeys map[string]string

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
		store:   s.store,
	}, true
}
