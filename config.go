package tidio

import (
	"io"
	"io/ioutil"
)

func NewConfig() *Config {
	c := &Config{
		activeLoggers: make(map[interface{}]*LogPrinter),
	}
	c.SetOutput(ioutil.Discard)
	return c
}

// Config holds reference to loggers for various objects.
type Config struct {
	out io.Writer

	debug bool

	activeLoggers map[interface{}]*LogPrinter
}

// SetOutput
func (me *Config) SetOutput(w io.Writer) { me.out = w }

// Debug
func (me *Config) Debug() bool { return me.debug }
func (me *Config) SetDebug(v bool) {
	me.debug = v
}

// Unreg removes the previously registered item if any.
func (me *Config) Unreg(v interface{}) {
	delete(me.activeLoggers, v)
}

func (me *Config) Log(v interface{}) *LogPrinter {
	l, found := me.activeLoggers[v]
	if !found {
		return nolog
	}
	return l
}

// Register creates a LogPrinter for the given objects.
// The registered LogPrinter is later retrieved with Config.Log(v)
func (me *Config) Register(v ...interface{}) *LogPrinter {
	if len(v) == 0 {
		panic("missing values in MLog")
	}

	first := v[0]
	l := NewLogPrinter(me.out)
	me.activeLoggers[first] = l
	for _, other := range v[1:] {
		me.activeLoggers[other] = l
	}
	return l
}

// ----------------------------------------

// Wrappers for default config Conf

func Register(v ...interface{}) *LogPrinter {
	return Conf.Register(v...)
}
func Log(v interface{}) *LogPrinter { return Conf.Log(v) }
func Unreg(v interface{})           { Conf.Unreg(v) }

var Conf = NewConfig()
