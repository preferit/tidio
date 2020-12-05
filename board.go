package tidio

import (
	"io"
	"io/ioutil"
)

func NewBoard() *Board {
	return &Board{
		logDest:       ioutil.Discard,
		activeLoggers: make(map[interface{}]*LogPrinter),
	}
}

// Board holds reference to loggers for various objects.
type Board struct {
	logDest       io.Writer
	activeLoggers map[interface{}]*LogPrinter
}

// Register
func (me *Board) Register(v interface{}) error {
	me.activeLoggers[v] = NewLogPrinter(me.logDest)
	return nil
}

// Unreg removes the previously registered item if any.
func (me *Board) Unreg(v interface{}) {
	delete(me.activeLoggers, v)
}

func (me *Board) Log(v interface{}) *LogPrinter {
	l, found := me.activeLoggers[v]
	if !found {
		return nolog
	}
	return l
}

func (me *Board) RLog(v ...interface{}) *LogPrinter {
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

var reg = NewBoard()

func RLog(v ...interface{}) *LogPrinter {
	return reg.RLog(v...)
}
func Log(v interface{}) *LogPrinter { return reg.Log(v) }

func Register(v interface{}) { reg.Register(v) }
func Unreg(v interface{})    { reg.Unreg(v) }
