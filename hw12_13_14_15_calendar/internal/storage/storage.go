package storage

import (
	"context"
	"time"
)

type Events interface {
	Create(ctx context.Context, event Event) (int, error)
	Update(ctx context.Context, id int, change Event) error
	Delete(ctx context.Context, userID int, id int) error
	DeleteAll(ctx context.Context, userID int) error
	ListAll(ctx context.Context, userID int) ([]Event, error)
	ListDay(ctx context.Context, userID int, date time.Time) ([]Event, error)
	ListWeek(ctx context.Context, userID int, date time.Time) ([]Event, error)
	ListMonth(ctx context.Context, userID int, date time.Time) ([]Event, error)
}

type Base interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
}

type Storage interface {
	Base
	Events
}
