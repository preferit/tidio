package tidio

import "github.com/gregoryv/fox"

// Logger wraps a configurable logger. Empty logger is usable as mute.
type Logger struct {
	fox.Logger
}

func (me *Logger) SetLogger(l fox.Logger) { me.Logger = l }

// Log is a noop unless fox logger is configured.
func (me *Logger) Log(v ...interface{}) {
	if me.Logger == nil {
		return
	}
	me.Logger.Log(v...)
}
