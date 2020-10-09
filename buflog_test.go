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
