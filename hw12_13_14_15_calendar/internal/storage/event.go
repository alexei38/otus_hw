package storage

import (
	"errors"
	"time"
)

var (
	ErrNotExistsEvent  = errors.New("event not found")
	ErrFailUpdateEvent = errors.New("update event failed")
	ErrEndDateOut      = errors.New("end date out of rage")
	ErrAlreadyExists   = errors.New("already exists")
)

type Event struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Start       time.Time `db:"start"`
	Stop        time.Time `db:"stop"`
	Description string    `db:"description"`
	UserID      int       `db:"user_id"`
	BeforeSend  time.Time `db:"before_send"`
}

type Notification struct {
	EventID    int
	EventTitle string
	Datetime   time.Time
	UserID     int
}

type User struct {
	ID    int
	Name  string
	Email string
}
