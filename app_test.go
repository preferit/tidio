package tidio

import (
	"net/http"
	"testing"

	"github.com/gregoryv/wolf"
)

func TestApp_supports_help_flag(t *testing.T) {
	cmd := wolf.NewTCmd("x", "-h")
	defer cmd.Cleanup()
	if NewApp(cmd).Run() != 0 {
		t.Error(cmd.Dump())
	}
}

func TestApp_has_sane_default_values(t *testing.T) {
	cmd := wolf.NewTCmd()
	defer cmd.Cleanup()
	app := NewApp(cmd)
	var started bool
	app.ListenAndServe = func(string, http.Handler) error {
		started = true
		return nil
	}
	app.Run()
	if !started {
		t.Error("didn't start:", cmd.Dump())
	}
}

func TestApp_fails_on_wrong_option(t *testing.T) {
	cmd := wolf.NewTCmd("x", "-wrong-option", "value")
	defer cmd.Cleanup()
	if NewApp(cmd).Run() == 0 {
		t.Error(cmd.Err)
	}
}

func Test_can_specify_state_file(t *testing.T) {
	cmd := wolf.NewTCmd("x", "-state", "other.file")
	defer cmd.Cleanup()
	app := NewApp(cmd)
	app.ListenAndServe = noopListenAndServe
	//app.SetLogger(t) fails as reload is done in the background
	if app.Run() != 0 {
		t.Error("failed")
	}
	// todo test reload without auto reload
	if app.Run() != 0 {
		t.Error("reload")
	}
}

func Test_state_file_cannot_be_written(t *testing.T) {
	cmd := wolf.NewTCmd("x", "-state", "/var/no-such")
	defer cmd.Cleanup()

	app := NewApp(cmd)
	app.ListenAndServe = noopListenAndServe

	if app.Run() == 0 {
		t.Fail()
	}
}

func Test_state_option_is_empty(t *testing.T) {
	cmd := wolf.NewTCmd("x", "-state", "")
	defer cmd.Cleanup()

	app := NewApp(cmd)
	app.ListenAndServe = noopListenAndServe

	if app.Run() == 0 {
		t.Error(cmd.Dump())
	}
}

var noopListenAndServe = func(string, http.Handler) error { return nil }
