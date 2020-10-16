package tidio

import "testing"

func TestSettings(t *testing.T) {
	setFail(t, UseHost(""), nil)
	setFail(t, InitialAccount{}, nil)
	setFail(t, UseLogger{}, nil)
	setFail(t, Credentials{}, nil)
}

func setFail(t *testing.T, s Setting, v interface{}) {
	t.Helper()
	if err := s.Set(v); err == nil {
		t.Error("should fail")
	}

}
