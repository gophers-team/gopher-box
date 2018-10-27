package main

import (
	"time"

	"github.com/jmoiron/sqlx"
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

func InitDb(dbFile string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbFile)
	return db, err
}
