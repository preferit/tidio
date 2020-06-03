package tidio

import "testing"

func Test_service(t *testing.T) {
	store := &Store{}
	apikeys := map[string]string{}
	service := NewService(apikeys, store)
	if service == nil {
		t.Fail()
	}
}
