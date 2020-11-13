package tidio

import "github.com/gregoryv/fox"

// Logger wraps a configurable logger. Empty logger is usable as mute.
type OptionalLogger struct {
	fox.Logger
}

func (me *OptionalLogger) SetLogger(l fox.Logger) { me.Logger = l }

// Log is a noop unless fox logger is configured.
func (me *OptionalLogger) Log(v ...interface{}) {
	if me.Logger == nil {
		return
	}
	me.Logger.Log(v...)
}
