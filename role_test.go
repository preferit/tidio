package tidio

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_role(t *testing.T) {
	store, cleanup := newTempStore(t)
	defer cleanup()
	role := &Role{
		account: "john",
		store:   store,
	}

	if got := role.Account(); got != "john" {
		t.Errorf("Account() returned %q", got)
	}

	t.Run("CreateTimesheet", func(t *testing.T) {
		format := "xx.txt"
		if err := role.CreateTimesheet(format, "john", aFile("x")); err == nil {
			t.Errorf("CreateTimesheet ok with fileformat %q", format)
		}
		if err := role.CreateTimesheet("202001.timesheet", "john", aFile(".")); err != nil {
			t.Error("failed to create timesheet", err)
		}
		if err := role.CreateTimesheet("202001.timesheet", "not-user", aFile(".")); err == nil {
			t.Error("created timesheet for not-user", err)
		}
	})

	t.Run("ReadTimesheet", func(t *testing.T) {
		filename := "199902.timesheet"
		role.CreateTimesheet(filename, "john", aFile("..."))
		if err := role.ReadTimesheet(ioutil.Discard, filename, "john"); err != nil {
			t.Error(err)
		}
	})
}

func aFile(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
