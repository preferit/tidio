package tidio

import (
	"io/ioutil"
	"strings"
	"testing"
)

func Test_role(t *testing.T) {
	service, cleanup := newTestService(t)
	defer cleanup()
	role := &Role{
		account: "john",
		service: service,
	}

	if got := role.Account(); got != "john" {
		t.Errorf("Account() returned %q", got)
	}
	t.Run("CreateTimesheet", func(t *testing.T) {
		file := ioutil.NopCloser(strings.NewReader("x"))
		format := "xx.txt"
		if err := role.CreateTimesheet(format, "john", file); err == nil {
			t.Errorf("CreateTimesheet ok with fileformat %q", format)
		}
		if err := role.CreateTimesheet("202001.timesheet", "john", file); err != nil {
			t.Error("failed to create timesheet", err)
		}
		if err := role.CreateTimesheet("202001.timesheet", "not-user", file); err == nil {
			t.Error("created timesheet for not-user", err)
		}
	})
}
