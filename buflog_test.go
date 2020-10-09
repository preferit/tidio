package tidio

import (
	"bytes"
	"fmt"
)

type BufferedLogger struct {
	bytes.Buffer
}

// Method
func (me *BufferedLogger) Log(v ...interface{}) {
	me.WriteString(fmt.Sprintln(v...))
}

func Buflog(srv *Service) *BufferedLogger {
	var buflog BufferedLogger
	srv.SetLogger(&buflog)
	return &buflog
}

// String
func (me *BufferedLogger) String() string {
	return "\n" + me.Buffer.String()
}
