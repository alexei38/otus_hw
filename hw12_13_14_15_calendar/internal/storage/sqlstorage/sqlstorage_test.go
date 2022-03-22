package sqlstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) { // nolint: funlen
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := New()
	s.db = sqlx.NewDb(db, "sqlmock")

	err = db.Ping()
	require.NoError(t, err)

	location, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)
	currDate := time.Date(2022, 0o3, 17, 15, 0o0, 0o0, 0o0, location)

	userID := 7
	var events []storage.Event
	for i := 0; i < 10; i++ {
		event := storage.Event{
			ID:          i,
			Title:       fmt.Sprintf("Title-%d", i),
			Start:       currDate,
			Stop:        currDate.Add(time.Hour),
			Description: fmt.Sprintf("Description-%d", i),
			UserID:      userID,
			BeforeSend:  time.Second * 10,
		}
		events = append(events, event)
	}

	t.Run("Create event", func(t *testing.T) {
		event := storage.Event{
			Title:       "Title",
			Start:       currDate.Add(time.Hour),
			Stop:        currDate.Add(time.Hour * 2),
			Description: "Description",
			UserID:      1,
			BeforeSend:  time.Second * 10,
		}
		query := `INSERT INTO events \(title, start, stop, description, user_id, before_send\)`
		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs(
			event.Title,
			event.Start,
			event.Stop,
			event.Description,
			event.UserID,
			event.BeforeSend,
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

		id, err := s.Create(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, 1)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Update event", func(t *testing.T) {
		change := storage.Event{
			ID:          1,
			Title:       "Title",
			Start:       currDate.Add(time.Hour * 2),
			Stop:        currDate.Add(time.Hour * 3),
			Description: "Description",
			UserID:      1,
			BeforeSend:  time.Second * 10,
		}
		query := `UPDATE events
			  SET title=.+,
			      start=.+,
			      stop=.+,
			      description=.+,
			      user_id=.+,
			      before_send=.+
		      WHERE id=.+`
		mock.ExpectPrepare(query).ExpectExec().WithArgs(
			change.Title,
			change.Start,
			change.Stop,
			change.Description,
			change.UserID,
			change.BeforeSend,
			change.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.Update(context.Background(), change.ID, change)
		require.NoError(t, err)

		// we make sure that all expectations were met
		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("Delete event", func(t *testing.T) {
		query := `DELETE FROM events WHERE id=.+ AND user_id=.+`
		mock.ExpectExec(query).WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.Delete(context.Background(), 1, 1)
		require.NoError(t, err)

		// we make sure that all expectations were met
		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("DeleteAll events", func(t *testing.T) {
		query := `DELETE FROM events WHERE user_id=.+`
		mock.ExpectExec(query).WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.DeleteAll(context.Background(), 1)
		require.NoError(t, err)

		// we make sure that all expectations were met
		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("ListAll events", func(t *testing.T) {
		var rows []*sqlmock.Rows
		for _, ev := range events {
			row := sqlmock.NewRows(
				[]string{"id", "title", "start", "stop", "description", "user_id", "before_send"},
			).AddRow(
				ev.ID,
				ev.Title,
				ev.Start,
				ev.Stop,
				ev.Description,
				ev.UserID,
				ev.BeforeSend,
			)
			rows = append(rows, row)
		}
		mock.ExpectQuery(
			`SELECT id, title, start, stop, description, user_id,
			before_send FROM events WHERE user_id=.+ ORDER BY start ASC`,
		).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows...)

		_, err = s.ListAll(context.Background(), userID)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("ListDay events", func(t *testing.T) {
		var rows []*sqlmock.Rows
		for _, ev := range events {
			row := sqlmock.NewRows(
				[]string{"id", "title", "start", "stop", "description", "user_id", "before_send"},
			).AddRow(
				ev.ID,
				ev.Title,
				ev.Start,
				ev.Stop,
				ev.Description,
				ev.UserID,
				ev.BeforeSend,
			)
			rows = append(rows, row)
		}
		findDate := currDate
		query := `SELECT id, title, start, stop, description, user_id, before_send
		FROM events
		WHERE user_id=.+ AND start::TIMESTAMP::DATE=.+ ORDER BY start ASC`
		mock.ExpectQuery(query).WithArgs(userID, findDate.Format("2006-01-02")).WillReturnRows(rows...)

		_, err = s.ListDay(context.Background(), userID, findDate)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("ListWeek events", func(t *testing.T) {
		var rows []*sqlmock.Rows
		for _, ev := range events {
			row := sqlmock.NewRows(
				[]string{"id", "title", "start", "stop", "description", "user_id", "before_send"},
			).AddRow(
				ev.ID,
				ev.Title,
				ev.Start,
				ev.Stop,
				ev.Description,
				ev.UserID,
				ev.BeforeSend,
			)
			rows = append(rows, row)
		}
		findDate := currDate
		year, week := findDate.ISOWeek()
		query := `SELECT id, title, start, stop, description, user_id, before_send
		FROM events
		WHERE user_id=.+
			AND extract\(isoyear from start\)=.+
			AND extract\(week from start\)=.+
		ORDER BY start ASC`
		mock.ExpectQuery(query).WithArgs(userID, year, week).WillReturnRows(rows...)

		_, err = s.ListWeek(context.Background(), userID, findDate)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("ListMonth events", func(t *testing.T) {
		var rows []*sqlmock.Rows
		for _, ev := range events {
			row := sqlmock.NewRows(
				[]string{"id", "title", "start", "stop", "description", "user_id", "before_send"},
			).AddRow(
				ev.ID,
				ev.Title,
				ev.Start,
				ev.Stop,
				ev.Description,
				ev.UserID,
				ev.BeforeSend,
			)
			rows = append(rows, row)
		}
		findDate := currDate
		year, month, _ := findDate.Date()
		query := `SELECT id, title, start, stop, description, user_id, before_send
		FROM events
		WHERE user_id=.+
		    AND extract\(isoyear from start\)=.+
		    AND extract\(month from start\)=.+
		ORDER BY start ASC`
		mock.ExpectQuery(query).WithArgs(userID, year, month).WillReturnRows(rows...)

		_, err = s.ListMonth(context.Background(), userID, findDate)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}
