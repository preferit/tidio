package tidio

import (
	"io"
	"net/http"
)

type API struct{}

// CreateTimesheet
func (me API) CreateTimesheet(loc string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest("POST", loc, body)
	if err != nil {
		return nil, err
	}
	return r, nil
}
