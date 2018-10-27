package main

import (
	"github.com/gophers-team/gopher-box/api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var DbFile = "/var/lib/gopher-box/events.db"

type EventType uint8

const (
	Heartbeat EventType = iota
	PillsTaken
)

type Schema struct {
	create string
	drop   string
}

var defaultSchema = Schema{
	create: `
CREATE TABLE events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	device_id	VARCHAR(127) NOT NULL,
    event_type  UNSIGNED INT8 NOT NULL,
    description VARCHAR(255),
    timestamp   DATETIME,
    created_at  DATETIME
);
`,
	drop: `
DROP TABLE events;
`,
}

type Event struct {
	Id          int           `db:"id"`
	DeviceId	api.DeviceID
	Type        EventType     `db:"event_type"`
	Description string        `db:"description"`
	Timestamp   time.Time     `db:"timestamp"`
	CreatedAt   time.Time     `db:"created_at"`
}

func InitDb() (*sqlx.DB, error) {
 	db, err := sqlx.Connect("sqlite3", DbFile)
 	return db, err
}
