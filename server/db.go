package main

import (
	"time"

	"github.com/gophers-team/gopher-box/api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type EventType uint8

const (
	Heartbeat EventType = iota
	PillsDispensed
)

type Schema struct {
	create string
	drop   string
}

type Event struct {
	Id          int `db:"id"`
	DeviceId    api.DeviceID
	Type        EventType `db:"event_type"`
	Description string    `db:"description"`
	Timestamp   time.Time `db:"timestamp"`
	CreatedAt   time.Time `db:"created_at"`
}

func InitDb(dbFile string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbFile)
	return db, err
}
