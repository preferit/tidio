package tidio

import "testing"

func TestCredentials_Set_on_nil(t *testing.T) {
	c := &Credentials{}
	if err := c.Set(nil); err == nil {
		t.Error("should fail")
	}
}
