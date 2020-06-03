package tidio

import "testing"

func Test_service_write_operations(t *testing.T) {
	_, cleanup := newTestService(t)
	defer cleanup()
}

func Test_service(t *testing.T) {
	service, cleanup := newTestService(t)
	defer cleanup()
	if service == nil {
		t.Fail()
	}
	if _, ok := service.IsAuthenticated("KEY"); !ok {
		t.Error("KEY is in apikeys")
	}
	if _, ok := service.IsAuthenticated(""); ok {
		t.Error("empty key ok")
	}
	if _, ok := service.IsAuthenticated("not there"); ok {
		t.Error("wrong key ok")
	}
}

func newTestService(t *testing.T) (*Service, func()) {
	store := &Store{}
	apikeys := APIKeys{
		"KEY": "john",
	}
	service := NewService(apikeys, store)
	cleanup := func() {}
	return service, cleanup
}

func Test_service_options(t *testing.T) {
	defer catchPanic(t)
	NewService(1)
}

func catchPanic(t *testing.T) {
	e := recover()
	if e == nil {
		t.Helper()
		t.Error("didn't panic")
	}
}
