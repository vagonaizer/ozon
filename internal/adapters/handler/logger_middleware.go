package handler

import (
	"net/http"
	"time"

	"github.com/vagonaizer/go-cart/internal/logger"
)

func LoggingMiddleware(next http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		logger.LogIncomingRequest(r)

		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		executionTime := time.Since(startTime)
		logger.LogOutgoingResponse(rw.statusCode, http.StatusText(rw.statusCode), executionTime)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
