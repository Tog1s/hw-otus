package storage

import (
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	Create(event *Event) (*Event, error)
	Update(event *Event) (*Event, error)
	Delete(event *Event) error
	DayEventList()
	WeekEventList()
	MonthEventList()
}

type Event struct {
	ID           uuid.UUID
	Title        string
	DateTime     time.Time
	EndTime      time.Time
	Description  string
	UserID       int
	NotifyBefore time.Time
}

type Notification struct {
	ID         uuid.UUID
	Title      string
	Date       time.Time
	SendToUser int
}
