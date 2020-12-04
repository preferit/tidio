package tidio

import (
	"net/http"
	"testing"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/wolf"
)

func TestApp_supports_help_flag(t *testing.T) {
	cmd := wolf.NewTCmd("x", "-h")
	defer cmd.Cleanup()
	if NewApp().Run(cmd) != 0 {
		t.Error(cmd.Dump())
	}
}

func TestApp_has_sane_default_values(t *testing.T) {
	app := NewApp(fox.Logging{t})
	var started bool
	app.ListenAndServe = func(string, http.Handler) error {
		started = true
		return nil
	}
	cmd := wolf.NewTCmd()
	defer cmd.Cleanup()
	app.Run(cmd)
	if !started {
		t.Error("didn't start:", cmd.Dump())
	}
}

func TestApp_fails_on_wrong_option(t *testing.T) {
	cmd := wolf.NewTCmd("x", "-wrong-option", "value")
	defer cmd.Cleanup()
	app := NewApp()
	app.SetLogger(t)
	code := app.Run(cmd)
	if code == 0 {
		t.Error(cmd.Err.String())
	}
}

func Test_can_specify_state_file(t *testing.T) {
	app := NewApp()
	app.ListenAndServe = noopListenAndServe
	cmd := wolf.NewTCmd("x", "-state", "other.file")
	defer cmd.Cleanup()
	if app.Run(cmd) != 0 {
		t.Error("failed")
	}

	if app.Run(cmd) != 0 {
		t.Error("should reload")
	}
}

func Test_state_file_cannot_be_written(t *testing.T) {
	app := NewApp()
	app.SetLogger(t)
	app.ListenAndServe = noopListenAndServe

	cmd := wolf.NewTCmd("x", "-state", "/var/no-such")
	defer cmd.Cleanup()
	if app.Run(cmd) == 0 {
		t.Error(cmd.Out.String(), cmd.Err.String())
	}
}

var noopListenAndServe = func(string, http.Handler) error { return nil }
