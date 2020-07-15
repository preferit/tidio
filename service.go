package tidio

import (
	"io/ioutil"
	"net/http"

	"github.com/gregoryv/fox"
)

func NewService() *Service {
	s := &Service{
		warn: fox.NewSyncLog(ioutil.Discard).Log,
	}
	return s
}

type Service struct {
	warn func(...interface{})
}

func (me *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
