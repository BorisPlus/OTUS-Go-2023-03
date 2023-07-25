package internalhttp

import (
	"net/http"
	"time"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)

func middleware(wrappedHandler http.Handler, logger interfaces.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartAt := time.Now()
		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)
		a := struct {
			StatusCode      int
			UserAgent       string
			ClientIPAddress string
			HTTPMethod      string
			HTTPVersion     string
			StartAt         time.Time
			Latency         time.Duration
		}{
			StatusCode:      lrw.StatusCode,
			UserAgent:       r.UserAgent(),
			ClientIPAddress: r.RemoteAddr,
			HTTPMethod:      r.Method,
			HTTPVersion:     r.Proto,
			StartAt:         StartAt,
			Latency:         time.Since(StartAt),
		}
		logger.Info("%+v", a)
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(writer http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{writer, 0}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
