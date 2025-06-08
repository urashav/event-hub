package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger логирует информацию о запросе
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем ResponseWriter, который может отслеживать статус
		wrw := newResponseWriter(w)

		next.ServeHTTP(wrw, r)

		// Логируем информацию о запросе
		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.RequestURI,
			wrw.status,
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
