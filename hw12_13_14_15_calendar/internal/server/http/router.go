package internalhttp

import "net/http"

func newRouter(app Application, logger Logger) *http.ServeMux {
	handler := Handler{App: app}
	router := http.NewServeMux()
	router.HandleFunc("/", loggingMiddleware(http.HandlerFunc(handler.index), logger))
	router.HandleFunc("/create", loggingMiddleware(http.HandlerFunc(handler.Create), logger))
	router.HandleFunc("/update", loggingMiddleware(http.HandlerFunc(handler.Update), logger))
	router.HandleFunc("/delete", loggingMiddleware(http.HandlerFunc(handler.Delete), logger))
	router.HandleFunc("/dayeventlist", loggingMiddleware(http.HandlerFunc(handler.DayEventList), logger))
	router.HandleFunc("/weekeventlist", loggingMiddleware(http.HandlerFunc(handler.WeekEventList), logger))
	router.HandleFunc("/montheventlist", loggingMiddleware(http.HandlerFunc(handler.MonthEventList), logger))
	return router
}
