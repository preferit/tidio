package tidio

func NewService(opts ...interface{}) *Service {
	service := &Service{}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case *Store:
			service.store = opt
		case map[string]string:
			service.apikeys = opt
		}
	}
	return service
}

type Service struct {
	store   *Store
	apikeys map[string]string
}
