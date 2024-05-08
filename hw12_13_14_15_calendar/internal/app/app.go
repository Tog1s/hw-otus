package app

import (
	"context"
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
	Create(*storage.Event) error
	Update(*storage.Event) error
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

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// if err := a.storage.Create(&storage.Event{ID: id, Title: tiitle}); err != nil {
	// 	return err
	// }
	return nil
}
