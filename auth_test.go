package tidio

import (
	"net/http"
	"testing"
)

func TestBasicAuth_Set(t *testing.T) {
	t.Run("on nil", func(t *testing.T) {
		c := NewBasicAuth(&Credentials{})
		if err := c.Set(nil); err == nil {
			t.Error("should fail")
		}
	})

	t.Run("on *http.Request", func(t *testing.T) {
		c := NewBasicAuth(&Credentials{})
		r, _ := http.NewRequest("GET", "http://example.com", nil)
		if err := c.Set(r); err != nil {
			t.Error("should work:", err)
		}
	})
}
