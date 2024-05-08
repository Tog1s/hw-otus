package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logger Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTime := time.Now()
		next.ServeHTTP(w, r)
		logger.Info(fmt.Sprintf(
			"%s %s %s %d %v %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			time.Since(requestTime),
			r.UserAgent(),
		))
	})
}
