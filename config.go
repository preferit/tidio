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
	out           io.Writer
	activeLoggers map[interface{}]*LogPrinter
}

// SetOutput
func (me *Config) SetOutput(w io.Writer) { me.out = w }

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
	l := NewLogPrinter(me.out)
	me.activeLoggers[first] = l
	for _, other := range v[1:] {
		me.activeLoggers[other] = l
	}
	return l
}

// ----------------------------------------

var Conf = NewConfig()

func RLog(v ...interface{}) *LogPrinter { return Conf.RLog(v...) }
func Log(v interface{}) *LogPrinter     { return Conf.Log(v) }

func Unreg(v interface{}) { Conf.Unreg(v) }
