package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		latency := time.Since(start)
		timestamp := time.Now().Format("02/Jan/2006:15:04:05 -0700")

		msg := fmt.Sprintf("%s [%s] %s %s %s %d %d %s",
			r.RemoteAddr,
			timestamp,
			r.Method,
			r.URL,
			r.Proto,
			lrw.statusCode,
			latency.Milliseconds(),
			r.UserAgent(),
		)

		s.logger.Info(msg)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	return lrw.ResponseWriter.Write(b)
}
