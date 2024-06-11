package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventAlreadyExist = errors.New("event already exist")
	ErrEventNotFound     = errors.New("event not found")
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]*storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]*storage.Event),
	}
}

func (s *Storage) Create(e *storage.Event) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.ID]; ok {
		return nil, ErrEventAlreadyExist
	}

	s.events[e.ID] = e
	return e, nil
}

func (s *Storage) Update(e *storage.Event) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] == nil {
		return nil, ErrEventNotFound
	}

	s.events[e.ID] = e
	return e, nil
}

func (s *Storage) Delete(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] == nil {
		return ErrEventNotFound
	}

	delete(s.events, e.ID)
	return nil
}

func (s *Storage) DayEventList(day time.Time) (map[uuid.UUID]*storage.Event, error) {
	end := day.Add(time.Hour * 24)
	return s.list(day, end)
}

func (s *Storage) WeekEventList(day time.Time) (map[uuid.UUID]*storage.Event, error) {
	end := day.Add(time.Hour * 24 * 7)
	return s.list(day, end)
}

func (s *Storage) MonthEventList(day time.Time) (map[uuid.UUID]*storage.Event, error) {
	end := day.Add(time.Hour * 24 * 7 * 30)
	return s.list(day, end)
}

func (s *Storage) list(start, end time.Time) (map[uuid.UUID]*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := make(map[uuid.UUID]*storage.Event)
	for _, e := range s.events {
		if e.DateTime.After(start) && e.EndTime.Before(end) {
			events[e.ID] = e
		}
	}
	return events, nil
}
