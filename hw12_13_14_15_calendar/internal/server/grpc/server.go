package internalgrpc

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Application interface {
	CreateEvent(*storage.Event) (*storage.Event, error)
	UpdateEvent(*storage.Event) (*storage.Event, error)
	DeleteEvent(*storage.Event) error
	DayEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
	WeekEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
	MonthEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
}

type Logger interface {
	Info(string)
	Error(string)
	Warn(string)
	Debug(string)
}

type Server struct {
	host   string
	port   string
	logger Logger
	app    Application
	server *grpc.Server
}

func New(logger Logger, app Application, host, port string) *Server {
	service := NewService(app, logger)
	server := grpc.NewServer(grpc.UnaryInterceptor(loggingMiddleware(logger)))
	pb.RegisterEventServiceServer(server, service)
	reflection.Register(server)
	s := &Server{
		host:   host,
		port:   port,
		logger: logger,
		app:    app,
		server: server,
	}
	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		fmt.Println("unable to create listener", err)
		panic(err)
	}

	s.logger.Info(fmt.Sprintf("GRPC-server %s:%s starting", s.host, s.port))
	if err = s.server.Serve(listener); err != nil {
		fmt.Println("unable to start grpc-server", err)
		panic(err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.server.GracefulStop()
	s.logger.Info("stopping GRPC-server")
	return nil
}
