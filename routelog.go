package tidio

import (
	"net/http"
	"time"
)

// NewRouteLog returns RouteLog using the package logger.
func NewRouteLog(dst *LogPrinter) *RouteLog {
	return &RouteLog{LogPrinter: dst}
}

type RouteLog struct {
	*LogPrinter
}

// Middleware returns a middleware that logs request and response status.
func (me *RouteLog) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{w, 200}
		next.ServeHTTP(&rec, r)
		me.Info(r.Method, r.URL, rec.status, time.Since(start))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
	rec.status = code
}
