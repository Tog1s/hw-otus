package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage"
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
	CreateEvent(*storage.Event) (*storage.Event, error)
	UpdateEvent(*storage.Event) (*storage.Event, error)
	DeleteEvent(*storage.Event) error
	DayEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
	WeekEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
	MonthEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	s := &Server{
		logger: logger,
		host:   host,
		port:   port,
		app:    app,
	}
	router := newRouter(s.app, s.logger)
	s.server = http.Server{
		Addr:              net.JoinHostPort(s.host, s.port),
		Handler:           router,
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
