package app

import (
	"context"

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

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
