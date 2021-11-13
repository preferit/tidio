package tidio

import "net/http"

func NewTracer(dst *LogPrinter) *Tracer {
	return &Tracer{dst: dst}
}

type Tracer struct {
	dst *LogPrinter
}

func (me *Tracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := Register(r).Buf()
		log.lgr.SetPrefix("TRACE: ")
		next.ServeHTTP(w, r)
		if log.Failed() {
			me.dst.Info(log.FlushString())
		}
		Unreg(r)
	})
}
