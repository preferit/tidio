package tidio

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gregoryv/fox"
)

func newTrace(dst fox.Logger, r *http.Request) (trace fox.Logger, cleanup func()) {
	trace, cl := NewTrace(dst)
	return trace, func() {
		cl(r.Method, r.URL)
	}
}

func NewTrace(dst fox.Logger) (trace fox.Logger, cleanup func(...interface{})) {
	var buf bytes.Buffer
	l := log.New(&buf, "", log.Lshortfile)
	trace = fox.LoggerFunc(func(v ...interface{}) {
		l.Output(3, fmt.Sprint(v...))
	})
	cleanup = func(v ...interface{}) {
		if buf.Len() > 0 {
			dst.Log("------------------------------")
			dst.Log(v...)
			dst.Log(buf.String())
			dst.Log("------------------------------")
		}
	}
	return
}
