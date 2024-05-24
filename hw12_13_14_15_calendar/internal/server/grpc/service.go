package internalgrpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	app    Application
	logger Logger
	pb.UnimplementedEventServiceServer
}

func NewService(app Application, logger Logger) *Service {
	return &Service{
		app:    app,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	storageEvent, err := makeEvent(request.Event)
	if err != nil {
		return nil, err
	}

	event, err := s.app.CreateEvent(storageEvent)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResponse{Id: event.ID.String()}, nil
}

func (s *Service) Update(_ context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	storageEvent, err := makeEvent(request.Event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = s.app.UpdateEvent(storageEvent)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateResponse{}, nil
}

func (s *Service) Delete(_ context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	event, err := makeEvent(request.Event)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.app.DeleteEvent(event)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteResponse{}, nil
}

func (s *Service) DayEventList(_ context.Context, request *pb.GetListRequest) (*pb.EventsResponse, error) {
	events, _ := s.app.DayEventList(request.StartDate.AsTime())
	response := collectEvents(events)

	return &pb.EventsResponse{Events: response}, nil
}

func (s *Service) WeekEventList(_ context.Context, request *pb.GetListRequest) (*pb.EventsResponse, error) {
	events, _ := s.app.WeekEventList(request.StartDate.AsTime())
	response := collectEvents(events)

	return &pb.EventsResponse{Events: response}, nil
}

func (s *Service) MonthEventList(_ context.Context, request *pb.GetListRequest) (*pb.EventsResponse, error) {
	events, _ := s.app.MonthEventList(request.StartDate.AsTime())
	response := collectEvents(events)

	return &pb.EventsResponse{Events: response}, nil
}

func collectEvents(events map[uuid.UUID]*storage.Event) []*pb.Event {
	collection := make([]*pb.Event, 0, len(events))
	for _, event := range events {
		collection = append(collection, makePbEvent(event))
	}

	return collection
}

func makeEvent(e *pb.Event) (*storage.Event, error) {
	id, err := uuid.Parse(e.Id)
	if e.Id != "" && err != nil {
		return nil, fmt.Errorf("unable to parse event id from request: %w", err)
	}

	userId, err := strconv.Atoi(e.UserId)
	if err != nil {
		return nil, fmt.Errorf("unable to parse user id from request: %w", err)
	}

	return &storage.Event{
			ID:           id,
			Title:        e.Title,
			DateTime:     e.Datetime.AsTime(),
			EndTime:      e.Endtime.AsTime(),
			Description:  e.Description,
			UserID:       userId,
			NotifyBefore: e.NotifyBefore.AsTime(),
		},
		nil
}

func makePbEvent(e *storage.Event) *pb.Event {
	userId := strconv.Itoa(e.UserID)
	return &pb.Event{
		Id:           e.ID.String(),
		Title:        e.Title,
		Datetime:     timestamppb.New(e.DateTime),
		Endtime:      timestamppb.New(e.EndTime),
		Description:  e.Description,
		UserId:       userId,
		NotifyBefore: timestamppb.New(e.NotifyBefore),
	}
}
