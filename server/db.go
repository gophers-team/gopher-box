package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type EventType uint8

const (
	Heartbeat EventType = iota
	PillsDispensed
)

type DispensingPlan struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type DispensingSchedule struct {
	Id               int       `db:"id"`
	PlanId           int       `db:"plan_id"`
	PillId           int       `db:"pill_id"`
	Amount           int       `db:"amount"`
	DispenseDow      int       `db:"dispense_dow"`
	ScheduleDuration int       `db:"schedule_duration"`
	DispenseTime     time.Time `db:"dispense_time"`
	CreatedAt        time.Time `db:"created_at"`
}

type DeviceDispensing struct {
	Id             int       `db:"id"`
	DeviceId       int       `db:"device_id"`
	ScheduleId     int       `db:"schedule_id"`
	PillsDispensed int       `db:"pills_dispensed"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
}

func InitDb(dbFile string, devel bool) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	if devel {
		fmt.Println("Using sqlite3 db")
		db, err = sqlx.Connect("sqlite3", dbFile)
	} else {
		fmt.Println("Using postgres db")
		db, err = sqlx.Connect(
			"postgres",
			"host=127.0.0.1 port=5432 user=box password=box dbname=box sslmode=disable",
		)
		db.DB.SetMaxIdleConns(2)
		db.DB.SetMaxOpenConns(2)
	}
	return db, err
}
