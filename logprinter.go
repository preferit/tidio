package tidio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func NewLogPrinter(w io.Writer) *LogPrinter {
	return &LogPrinter{
		lgr: log.New(w, "", log.Lshortfile),
	}
}

var nolog = &LogPrinter{
	lgr: log.New(ioutil.Discard, "", 0),
}

type LogPrinter struct {
	buf    bytes.Buffer // if buffered
	lgr    *log.Logger
	writes int
}

// Buf makes the log printer buffered. Use Flush to get the contents.
func (me *LogPrinter) Buf() *LogPrinter {
	me.lgr.SetOutput(&me.buf)
	return me
}

// FlushString
func (me *LogPrinter) FlushString() string {
	return string(me.Flush())
}

// Flush returns the buffered bytes if any and resets the buffer.
func (me *LogPrinter) Flush() []byte {
	defer me.buf.Reset()
	return me.buf.Bytes()
}

// Info
func (me *LogPrinter) Info(v ...interface{}) {
	me.lgr.Output(2, fmt.Sprintln(v...))
	me.writes++
}

func (me *LogPrinter) Log(v ...interface{}) {
	me.lgr.Output(2, fmt.Sprintln(v...))
	me.writes++
}
