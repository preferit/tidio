package tidio

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func NewBoard() *Board {
	return &Board{
		activeLoggers: make(map[interface{}]*LogPrinter),
	}
}

// Register
func (me *Board) Register(v interface{}) {
	out := os.Stderr
	me.activeLoggers[v] = &LogPrinter{
		log: log.New(out, "", log.Lshortfile),
	}
}

// Unreg removes the previously registered item if any.
func (me *Board) Unreg(v interface{}) {
	delete(me.activeLoggers, v)
}

func (me *Board) Log(v interface{}) *LogPrinter {
	l, found := me.activeLoggers[v]
	if !found {
		l = &LogPrinter{
			log: log.New(ioutil.Discard, "", 0),
		}
		me.activeLoggers[v] = l
	}
	return l
}

type LogPrinter struct {
	buf    bytes.Buffer // if buffered
	log    *log.Logger
	writes int
}

// Buf makes the log printer buffered. Use Flush to get the contents.
func (me *LogPrinter) Buf() *LogPrinter {
	me.log = log.New(&me.buf, "", log.Lshortfile)
	return me
}

// Flush returns the buffered bytes if any and resets the buffer.
func (me *LogPrinter) Flush() []byte {
	defer me.buf.Reset()
	return me.buf.Bytes()
}

// Info
func (me *LogPrinter) Info(v ...interface{}) {
	me.log.Output(2, fmt.Sprint(v...))
	me.writes++
}

// ----------------------------------------

var reg = NewBoard()

func RLog(v interface{}) *LogPrinter {
	Register(v)
	return Log(v)
}
func Log(v interface{}) *LogPrinter { return reg.Log(v) }

func Register(v interface{}) { reg.Register(v) }
func Unreg(v interface{})    { reg.Unreg(v) }

// Board holds reference to loggers for various objects.
type Board struct {
	activeLoggers map[interface{}]*LogPrinter
}
