package memorystorage

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type events map[int]storage.Event

type Memory struct {
	events events
	mu     sync.RWMutex
	nextID int32
}

func New() *Memory {
	s := &Memory{}
	s.events = make(events)
	return s
}

func (s *Memory) Create(ctx context.Context, event storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.getID()
	if event.Start.After(event.Stop) {
		return 0, storage.ErrEndDateOut
	}

	if _, ok := s.events[id]; ok {
		return 0, storage.ErrAlreadyExists
	}
	s.events[id] = storage.Event{
		ID:          id,
		Title:       event.Title,
		Start:       event.Start,
		Stop:        event.Stop,
		Description: event.Description,
		UserID:      event.UserID,
		BeforeSend:  event.BeforeSend,
	}
	return id, nil
}

func (s *Memory) Update(ctx context.Context, id int, change storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	event, ok := s.events[id]
	if !ok {
		return storage.ErrNotExistsEvent
	}
	if id != change.ID {
		return storage.ErrFailUpdateEvent
	}
	if event.UserID != change.UserID {
		return storage.ErrNotExistsEvent
	}
	if change.Start.After(change.Stop) {
		return storage.ErrEndDateOut
	}
	s.events[id] = storage.Event{
		ID:          id,
		Title:       change.Title,
		Start:       change.Start,
		Stop:        change.Stop,
		Description: change.Description,
		UserID:      change.UserID,
		BeforeSend:  change.BeforeSend,
	}
	return nil
}

func (s *Memory) Delete(ctx context.Context, eventID int, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	event, ok := s.events[eventID]
	if !ok {
		return storage.ErrNotExistsEvent
	}
	if event.UserID != userID {
		return storage.ErrNotExistsEvent
	}
	delete(s.events, eventID)
	return nil
}

func (s *Memory) DeleteAll(ctx context.Context, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range s.events {
		if event.UserID == userID {
			delete(s.events, event.ID)
		}
	}
	return nil
}

func (s *Memory) Connect(ctx context.Context, _ string) error {
	return nil
}

func (s *Memory) Close(ctx context.Context) error {
	return nil
}

func (s *Memory) ListAll(ctx context.Context, userID int) ([]storage.Event, error) {
	var result []storage.Event
	for _, ev := range s.events {
		if ev.UserID == userID {
			result = append(result, ev)
		}
	}
	return result, nil
}

func (s *Memory) ListDay(ctx context.Context, userID int, date time.Time) ([]storage.Event, error) {
	var result []storage.Event
	year, month, day := date.Date()
	for _, ev := range s.events {
		evYear, evMonth, evDay := ev.Start.Date()
		if ev.UserID == userID && year == evYear && month == evMonth && day == evDay {
			result = append(result, ev)
		}
	}
	return result, nil
}

func (s *Memory) ListWeek(ctx context.Context, userID int, date time.Time) ([]storage.Event, error) {
	var result []storage.Event
	year, week := date.ISOWeek()
	for _, ev := range s.events {
		evYear, evWeek := ev.Start.ISOWeek()
		if ev.UserID == userID && year == evYear && week == evWeek {
			result = append(result, ev)
		}
	}
	return result, nil
}

func (s *Memory) ListMonth(ctx context.Context, userID int, date time.Time) ([]storage.Event, error) {
	var result []storage.Event
	year, month, _ := date.Date()
	for _, ev := range s.events {
		evYear, evMonth, _ := ev.Start.Date()
		if ev.UserID == userID && year == evYear && month == evMonth {
			result = append(result, ev)
		}
	}
	return result, nil
}

func (s *Memory) getID() int {
	currID := s.nextID
	atomic.AddInt32(&s.nextID, 1)
	return int(currID)
}
