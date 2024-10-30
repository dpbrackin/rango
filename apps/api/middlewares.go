package main

import (
	"context"
	"log"
	"net/http"
	"rango/auth"
	"rango/router"
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

// AuthMidleware checks for a loged in user and passes it into the context.
// If their is no logged in user, it will reject the request with a 401.
func AuthMiddleware(srv *auth.AuthService) router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID, err := r.Cookie("sessionID")

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}

			user, err := srv.AuthenticateSession(r.Context(), sessionID.Value)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
