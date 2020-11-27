package tidio

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gregoryv/fox"
)

func newTrace(dst fox.Logger, r *http.Request) (trace fox.Logger, cleanup func()) {
	var buf bytes.Buffer
	l := log.New(&buf, "", log.Lshortfile)
	trace = fox.LoggerFunc(func(v ...interface{}) {
		l.Output(3, fmt.Sprint(v...))
	})
	cleanup = func() {
		if buf.Len() > 0 {
			dst.Log("------------------------------")
			dst.Log(r.Method, r.URL)
			dst.Log(buf.String())
			dst.Log("------------------------------")
		}
	}
	return
}
