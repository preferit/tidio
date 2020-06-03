package tidio

import "testing"

func Test_role(t *testing.T) {
	role := &Role{
		account: "john",
	}
	if got := role.Account(); got != "john" {
		t.Errorf("Account() returned %q", got)
	}
}
