package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

const timeFormat = ""

type Storage struct {
	ctx context.Context
	db  *sqlx.DB
}

func New(ctx context.Context) *Storage {
	return &Storage{ctx: ctx}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		fmt.Println(err)
	}
	s.db = db
	return err
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (s *Storage) Create(e *storage.Event) error {
	query := `insert into events(id, title, date_time, end_time, description, user_id, notify_before)
		values($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.db.ExecContext(s.ctx, query, e.ID, e.Title, e.DateTime, e.EndTime, e.Description, e.UserID, e.NotifyBefore)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Update(e *storage.Event) error {
	query := `UPDATE events
		SET title = $1, date_time = $2, end_time = $3, description = $4, user_id = $5, notify_before = $6
		WHERE id = $7
	`

	_, err := s.db.ExecContext(s.ctx, query, e.Title, e.DateTime, e.EndTime, e.Description, e.UserID, e.NotifyBefore, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(e *storage.Event) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := s.db.ExecContext(s.ctx, query, e.ID)
	if err != nil {
		return err
	}
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
	query := `SELECT title, start_date, end_date, description, owner_id, remind_at 
			FROM events
			WHERE start_date>=$1 AND start_date<$2
			ORDER BY start_date
	`
	rows, err := s.db.QueryContext(s.ctx, query, start.Format(timeFormat), end.Format(timeFormat))
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]*storage.Event)
	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(&e.Title, &e.DateTime, &e.EndTime, &e.Description, &e.UserID, &e.NotifyBefore); err != nil {
			return nil, err
		}
		result[e.ID] = &e
	}

	return result, nil
}
