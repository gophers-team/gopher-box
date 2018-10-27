package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var DbFile = "/var/lib/gopher-box/events.db"

type EventType uint

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
    type        VARCHAR(127) NOT NULL,
    description VARCHAR(255),
    timestamp   INTEGER,
    created_at  INTEGER
);
`,
	drop: `
DROP TABLE events;
`,
}

type Event struct {
	Id          int           `db:"id"`
	Type        EventType     `db:"type"`
	Description string        `db:"description"`
	Timestamp   string        `db:"timestamp"`
	CreatedAt   time.Duration `db:"created_at"`
}

func InitDb() (*sqlx.DB, error) {
 	db, err := sqlx.Connect("sqlite3", DbFile)
 	return db, err
}
