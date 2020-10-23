package tidio

import (
	"testing"

	"github.com/gregoryv/ant"
)

func TestSettings(t *testing.T) {
	setFail(t, InitialAccount{}, nil)
	setFail(t, Logging{}, nil)
	setFail(t, ErrorHandling(t.Fatal), nil)
}

func setFail(t *testing.T, s ant.Setting, v interface{}) {
	t.Helper()
	if err := s.Set(v); err == nil {
		t.Error("should fail")
	}

}
