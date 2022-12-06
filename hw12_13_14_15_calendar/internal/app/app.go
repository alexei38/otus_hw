package app

import (
	"context"
	"fmt"
	"time"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	storage storage.Storage
}

type Logger interface { // TODO
}

func New(storage storage.Storage) *App {
	return &App{
		storage,
	}
}

func (a *App) Create(ctx context.Context, id int, title string, start time.Time, stop time.Time, description string, userID int, beforeSend time.Time) (int, error) {
	eventID, err := a.storage.Create(
		ctx, storage.Event{
			ID:          id,
			Title:       title,
			Start:       start,
			Stop:        stop,
			Description: description,
			UserID:      userID,
			BeforeSend:  beforeSend,
		},
	)
	fmt.Println(err)
	fmt.Println(eventID)
	all, _ := a.storage.ListAll(ctx, 0)
	fmt.Printf("%+v\n", all)
	if err != nil {
		return 0, err
	}
	return eventID, nil
}
