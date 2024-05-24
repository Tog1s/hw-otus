package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(string)
	Warn(string)
	Debug(string)
	Error(string)
}

type Storage interface {
	Create(*storage.Event) (*storage.Event, error)
	Update(*storage.Event) (*storage.Event, error)
	Delete(*storage.Event) error
	DayEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
	WeekEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
	MonthEventList(time.Time) (map[uuid.UUID]*storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(e *storage.Event) (*storage.Event, error) {
	event, err := a.storage.Create(e)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (a *App) UpdateEvent(e *storage.Event) (*storage.Event, error) {
	event, err := a.storage.Update(e)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (a *App) DeleteEvent(e *storage.Event) error {
	err := a.storage.Delete(e)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DayEventList(day time.Time) (map[uuid.UUID]*storage.Event, error) {
	events, err := a.storage.DayEventList(day)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (a *App) WeekEventList(day time.Time) (map[uuid.UUID]*storage.Event, error) {
	events, err := a.storage.WeekEventList(day)
	if err != nil {
		return nil, err
	}
	return events, nil

}

func (a *App) MonthEventList(day time.Time) (map[uuid.UUID]*storage.Event, error) {
	events, err := a.storage.MonthEventList(day)
	if err != nil {
		return nil, err
	}
	return events, nil
}
