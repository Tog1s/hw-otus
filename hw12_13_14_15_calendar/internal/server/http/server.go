package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	host   string
	port   string
	server http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Info(string)
	Error(string)
	Warn(string)
	Debug(string)
}

type Application interface {
	CreateEvent(context.Context, string, string) error
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	s := &Server{
		logger: logger,
		host:   host,
		port:   port,
		app:    app,
	}
	router := http.NewServeMux()
	router.HandleFunc("/", index)
	loggedRouter := loggingMiddleware(router, s.logger)
	s.server = http.Server{
		Addr:              net.JoinHostPort(s.host, s.port),
		Handler:           loggedRouter,
		ReadHeaderTimeout: time.Second * 90,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("HTTP-server %s:%s starting", s.host, s.port))
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.Shutdown(ctx)
	return nil
}

func index(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Index page\n"))
	w.WriteHeader(200)
}
