package main

import (
	"log"
	"net/http"
	"time"
)

// httpResponseWriter observes calls to another http.ResponseWriter that change
// the HTTP status code.
type httpResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w httpResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func Logger(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			status := http.StatusOK
			start := time.Now()

			addr := r.Header.Get("X-Real-IP")
			if addr == "" {
				addr = r.Header.Get("X-Forwarded-For")
				if addr == "" {
					addr = r.RemoteAddr
				}
			}
			rw := httpResponseWriter{w, status}

			logger.Printf("Started %s %s for %s", r.Method, r.URL.Path, addr)
			defer logger.Printf("Completed %v %s in %v\n", rw.status, http.StatusText(rw.status), time.Since(start))

			h.ServeHTTP(rw, r)
		})
	}
}
