package main

import (
	"log"
	"net/http"
	"time"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (writer *CustomResponseWriter) WriteHeader(code int) {
	writer.statusCode = code
	writer.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs the status code, path and response time of a requet
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqStart := time.Now()
		path := r.URL.Path

		responseWriter := &CustomResponseWriter{
			ResponseWriter: w,
			statusCode:     0,
		}

		next.ServeHTTP(responseWriter, r)

		log.Printf("[%d] %v %v", responseWriter.statusCode, path, time.Since(reqStart))
	})
}
