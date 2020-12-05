package tidio

import (
	"io"
	"io/ioutil"
)

func NewConfig() *Config {
	return &Config{
		logDest:       ioutil.Discard,
		activeLoggers: make(map[interface{}]*LogPrinter),
	}
}

// Config holds reference to loggers for various objects.
type Config struct {
	logDest       io.Writer
	activeLoggers map[interface{}]*LogPrinter
}

// Register
func (me *Config) Register(v interface{}) error {
	me.activeLoggers[v] = NewLogPrinter(me.logDest)
	return nil
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

func (me *Config) RLog(v ...interface{}) *LogPrinter {
	if len(v) == 0 {
		panic("missing values in MLog")
	}

	first := v[0]
	Register(first)
	l := Log(first)
	for _, other := range v[1:] {
		me.activeLoggers[other] = l
	}
	return l
}

// ----------------------------------------

var reg = NewConfig()

func RLog(v ...interface{}) *LogPrinter {
	return reg.RLog(v...)
}
func Log(v interface{}) *LogPrinter { return reg.Log(v) }

func Register(v interface{}) { reg.Register(v) }
func Unreg(v interface{})    { reg.Unreg(v) }
