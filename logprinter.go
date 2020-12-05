package tidio

import (
	"bytes"
	"fmt"
	"log"
)

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

// FlushString
func (me *LogPrinter) FlushString() string {
	return string(me.Flush())
}

// Info
func (me *LogPrinter) Info(v ...interface{}) {
	me.log.Output(2, fmt.Sprintln(v...))
	me.writes++
}

func (me *LogPrinter) Log(v ...interface{}) {
	me.log.Output(2, fmt.Sprintln(v...))
	me.writes++
}
