package tidio

import (
	"strings"
	"testing"
)

func Test_loggers(t *testing.T) {
	log := RLog(t).Buf()
	defer Unreg(t)

	somefunc(t) // should log
	reg.Unreg(t)
	somefunc(t) // no logger registered

	got := string(log.Flush())
	if strings.Count(got, "hello") != 1 {
		t.Errorf("cached log\n%s", got)
		t.Error("writes", log.writes)
	}
}

func somefunc(t *testing.T) {
	Log(t).Info("hello")
	Log(t).Info("world")
}
