package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func newEvent(userID int, start, stop time.Time) storage.Event {
	event := storage.Event{
		Title:       faker.Name(),
		Start:       start,
		Stop:        stop,
		Description: faker.Name(),
		UserID:      userID,
		BeforeSend:  time.Now().Add(time.Second * -10),
	}
	return event
}

func TestStorage(t *testing.T) { // nolint: funlen
	ctx := context.Background()
	t.Run("working memory storage", func(t *testing.T) {
		userID := 15
		s := New()
		require.Len(t, s.events, 0)
		event := newEvent(userID, time.Now(), time.Now().Add(time.Hour))
		id, err := s.Create(ctx, event)
		require.NoError(t, err)
		require.Len(t, s.events, 1)

		createdEvent := s.events[id]
		require.Equal(t, createdEvent.Title, event.Title)
		require.Equal(t, createdEvent.Start, event.Start)
		require.Equal(t, createdEvent.Description, event.Description)
		require.Equal(t, createdEvent.UserID, userID)
		require.Equal(t, createdEvent.BeforeSend, event.BeforeSend)

		err = s.Delete(ctx, createdEvent.ID, userID)
		require.NoError(t, err)
		require.Len(t, s.events, 0)
	})

	t.Run("delete all", func(t *testing.T) {
		userID := 20
		s := New()
		require.Len(t, s.events, 0)
		for i := 0; i < 10; i++ {
			_, err := s.Create(
				ctx,
				newEvent(userID, time.Now(), time.Now().Add(time.Hour)),
			)
			require.NoError(t, err)
		}
		require.Len(t, s.events, 10)
		err := s.DeleteAll(ctx, userID)
		if err != nil {
			require.NoError(t, err)
		}
		require.Len(t, s.events, 0)

		err = s.DeleteAll(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("update memory cache", func(t *testing.T) {
		userID := 50
		s := New()
		require.Len(t, s.events, 0)
		event := newEvent(
			userID,
			time.Now(),
			time.Now().Add(time.Hour),
		)
		id, err := s.Create(ctx, event)
		require.NoError(t, err)
		require.Len(t, s.events, 1)

		updatedEvent := newEvent(
			userID,
			time.Now().Add(time.Hour*2),
			time.Now().Add(time.Hour*5),
		)

		err = s.Update(ctx, id, updatedEvent)
		require.NoError(t, err)
		memEvent := s.events[id]
		require.Equal(t, updatedEvent.Title, memEvent.Title)
		require.Equal(t, updatedEvent.Start, memEvent.Start)
		require.Equal(t, updatedEvent.Description, memEvent.Description)
		require.Equal(t, updatedEvent.UserID, memEvent.UserID)
		require.Equal(t, updatedEvent.BeforeSend, memEvent.BeforeSend)

		err = s.Delete(ctx, updatedEvent.ID, userID)
		require.NoError(t, err)
		require.Len(t, s.events, 0)
	})

	t.Run("update event errors", func(t *testing.T) {
		user := 50
		user2 := 100
		s := New()
		require.Len(t, s.events, 0)
		event := newEvent(
			user,
			time.Now().Add(time.Hour*2),
			time.Now().Add(time.Hour*5),
		)
		id, err := s.Create(ctx, event)
		require.NoError(t, err)
		require.Len(t, s.events, 1)

		updatedEvent := newEvent(
			user2,
			time.Now().Add(time.Hour*2),
			time.Now().Add(time.Hour*5),
		)

		err = s.Update(ctx, id, updatedEvent)
		require.ErrorIs(t, err, storage.ErrNotExistsEvent)
	})

	t.Run("List events", func(t *testing.T) {
		location, err := time.LoadLocation("Europe/Moscow")
		require.NoError(t, err)

		currDate := time.Date(2022, 0o3, 17, 15, 0o0, 0o0, 0o0, location)
		prevYear := currDate.AddDate(-1, 0, 0)
		prevMonth := currDate.AddDate(0, -1, 0)
		prevWeek := currDate.AddDate(0, 0, -7)
		prevDay := currDate.AddDate(0, 0, -1)
		nextDay := currDate.AddDate(0, 0, 1)
		nextWeek := currDate.AddDate(0, 0, 7)
		nextMonth := currDate.AddDate(0, 1, 0)
		nextYear := currDate.AddDate(1, -1, 0)
		userID := 15
		s := New()
		require.Len(t, s.events, 0)

		dates := []time.Time{
			prevYear,
			prevMonth,
			prevWeek,
			prevDay,
			currDate,
			nextDay,
			nextWeek,
			nextMonth,
			nextYear,
		}
		for i, eventTime := range dates {
			_, err := s.Create(ctx, newEvent(
				userID,
				eventTime,
				eventTime.Add(time.Hour),
			))
			require.NoError(t, err)
			require.Len(t, s.events, i+1)
		}

		events, err := s.ListAll(ctx, userID)
		require.NoError(t, err)
		require.Len(t, events, len(dates))

		// ListDay
		events, err = s.ListDay(ctx, userID, currDate)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, events[0].Start, currDate)

		// ListWeek
		events, err = s.ListWeek(ctx, userID, currDate)
		require.NoError(t, err)
		require.Len(t, events, 3)
		var weekDates []time.Time
		for _, e := range events {
			weekDates = append(weekDates, e.Start)
		}
		require.Contains(t, weekDates, prevDay)
		require.Contains(t, weekDates, currDate)
		require.Contains(t, weekDates, nextDay)

		// ListMonth
		events, err = s.ListMonth(ctx, userID, currDate)
		require.NoError(t, err)
		require.Len(t, events, 5)
		var monthDates []time.Time
		for _, e := range events {
			monthDates = append(monthDates, e.Start)
		}
		require.Contains(t, monthDates, prevWeek)
		require.Contains(t, monthDates, prevDay)
		require.Contains(t, monthDates, currDate)
		require.Contains(t, monthDates, nextDay)
		require.Contains(t, monthDates, nextWeek)

		// Very old List
		oldDate := currDate.AddDate(-1000, 0, 0)
		events, err = s.ListDay(ctx, userID, oldDate)
		require.NoError(t, err)
		require.Len(t, events, 0)

		events, err = s.ListWeek(ctx, userID, oldDate)
		require.NoError(t, err)
		require.Len(t, events, 0)

		events, err = s.ListMonth(ctx, userID, oldDate)
		require.NoError(t, err)
		require.Len(t, events, 0)
	})

	t.Run("parallel create events", func(t *testing.T) {
		s := New()
		count := 5000
		wg := &sync.WaitGroup{}
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(userID int) {
				event := newEvent(userID, time.Now(), time.Now().Add(time.Hour))
				s.Create(ctx, event)
				wg.Done()
			}(i)
		}
		wg.Wait()
		require.Len(t, s.events, count)
	})
}
