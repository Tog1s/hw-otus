package internalhttp

import "net/http"

func newRouter(app Application, logger Logger) *http.ServeMux {
	handler := Handler{App: app}
	router := http.NewServeMux()
	router.HandleFunc("/", loggingMiddleware(http.HandlerFunc(handler.index), logger))
	router.HandleFunc("/create", loggingMiddleware(http.HandlerFunc(handler.Create), logger))

	return router
}
