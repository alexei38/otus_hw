package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	db *sqlx.DB
}

func (s *SQLStorage) ListAll(ctx context.Context, userID int) ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT id, title, start, stop, description, user_id, before_send 
		FROM events 
		WHERE user_id=$1 ORDER BY start ASC`
	err := s.db.SelectContext(ctx, &events, query, userID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) ListDay(ctx context.Context, userID int, date time.Time) ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT id, title, start, stop, description, user_id, before_send 
		FROM events 
		WHERE user_id=$1 AND start::TIMESTAMP::DATE=$2 ORDER BY start ASC`
	err := s.db.SelectContext(ctx, &events, query, userID, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) ListWeek(ctx context.Context, userID int, date time.Time) ([]storage.Event, error) {
	var events []storage.Event
	year, week := date.ISOWeek()
	query := `SELECT id, title, start, stop, description, user_id, before_send 
		FROM events 
		WHERE user_id=$1 
			AND extract(isoyear from start)=$2 
			AND extract(week from start)=$3 
		ORDER BY start ASC`
	err := s.db.SelectContext(ctx, &events, query, userID, year, week)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) ListMonth(ctx context.Context, userID int, date time.Time) ([]storage.Event, error) {
	var events []storage.Event
	year, month, _ := date.Date()
	query := `SELECT id, title, start, stop, description, user_id, before_send 
			  FROM events 
			  WHERE user_id=$1 
			       AND extract(isoyear from start)=$2 
			       AND extract(month from start)=$3 
			  ORDER BY start ASC`
	err := s.db.SelectContext(ctx, &events, query, userID, year, month)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) Connect(ctx context.Context, dsn string) (err error) {
	s.db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *SQLStorage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *SQLStorage) Create(ctx context.Context, event storage.Event) (int, error) {
	// https://github.com/jmoiron/sqlx/issues/154#issuecomment-148216948

	if event.Start.After(event.Stop) {
		return 0, storage.ErrEndDateOut
	}

	query := `INSERT INTO events (title, start, stop, description, user_id, before_send)
			  VALUES (:title, :start, :stop, :description, :user_id, :before_send) RETURNING id`

	stmt, err := s.db.PrepareNamedContext(
		ctx,
		query,
	)
	if err != nil {
		return 0, err
	}
	var id int
	err = stmt.GetContext(ctx, &id, event)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SQLStorage) Update(ctx context.Context, id int, change storage.Event) error {
	if change.Start.After(change.Stop) {
		return storage.ErrEndDateOut
	}

	change.ID = id
	query := `UPDATE events 
			  SET title=:title, 
			      start=:start, 
			      stop=:stop,
			      description=:description, 
			      user_id=:user_id, 
			      before_send=:before_send 
		      WHERE id=:id`
	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, change)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) Delete(ctx context.Context, id int, userID int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func (s *SQLStorage) DeleteAll(ctx context.Context, userID int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE user_id=$1", userID)
	return err
}

func New() *SQLStorage {
	return &SQLStorage{}
}
