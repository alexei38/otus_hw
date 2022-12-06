package http

import (
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		resp := &responseData{
			status: 0,
			size:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: rw,
			responseData:   resp,
		}
		h.ServeHTTP(&lrw, r)

		duration := time.Since(start)
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		log.WithFields(log.Fields{
			"addr":       host,
			"uri":        r.RequestURI,
			"proto":      r.Proto,
			"method":     r.Method,
			"duration":   duration,
			"status":     resp.status,
			"size":       resp.size,
			"user_agent": r.UserAgent(),
		}).Info("request completed")
	}
	return http.HandlerFunc(logFn)
}
