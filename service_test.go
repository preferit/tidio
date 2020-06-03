package tidio

import "testing"

func Test_service(t *testing.T) {
	store := &Store{}
	apikeys := APIKeys{}
	service := NewService(apikeys, store)
	if service == nil {
		t.Fail()
	}
}
