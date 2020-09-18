package tidio

import (
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/rs"
)

func TestService_AddAccount(t *testing.T) {
	srv := NewService()
	ok, bad := asserter.NewErrors(t)
	ok(srv.AddAccount("john", "secret"))
	bad(srv.AddAccount("john", "secret"))
	bad(srv.AddAccount("root", "secret"))
	_, err := rs.Root.Use(srv.sys).Stat("/api/timesheets/john")
	ok(err)
}

func TestService_RestoreState(t *testing.T) {
	srv := NewService()
	ok, bad := asserter.NewErrors(t)
	bad(srv.RestoreState("no-such-file"))
	ok(srv.RestoreState(""))

}
