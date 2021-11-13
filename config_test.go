package tidio

import (
	"strings"
	"testing"
)

func Test_loggers(t *testing.T) {
	log := Register(t).Buf()
	defer Unreg(t)

	somefunc(t) // should log
	Conf.Unreg(t)
	somefunc(t) // no logger registered

	got := log.FlushString()
	if strings.Count(got, "hello") != 1 {
		t.Errorf("cached log\n%s", got)
		t.Error("writes", log.writes)
	}
}

func Test_Register_panics(t *testing.T) {
	defer catchPanic(t)
	Register()
}

func somefunc(t *testing.T) {
	Log(t).Info("hello")
	Log(t).Info("world")
}
