package frontend

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
)

type logger struct {
	*zap.Logger
	next http.Handler
}

func newLogMiddleware(next http.Handler, l *zap.Logger) http.Handler {
	return &logger{
		Logger: l,
		next:   next,
	}
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := httpsnoop.CaptureMetrics(l.next, w, r)
	l.Info(
		"request complete",
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.Int("code", m.Code),
		zap.Duration("duration", m.Duration),
		zap.Int64("written", m.Written),
	)
}
