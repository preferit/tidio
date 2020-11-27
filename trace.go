package tidio

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gregoryv/fox"
)

func newTrace(dst fox.Logger, r *http.Request) (trace *log.Logger, cleanup func()) {
	var buf bytes.Buffer
	trace = log.New(&buf, "", log.Lshortfile)
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
