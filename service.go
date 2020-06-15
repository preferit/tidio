package tidio

import "net/http"

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
