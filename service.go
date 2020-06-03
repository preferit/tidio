package tidio

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
		}
	}
	return service
}

type Service struct {
	store   *Store
	apikeys APIKeys
}

type APIKeys map[string]string
