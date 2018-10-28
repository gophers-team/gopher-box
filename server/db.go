package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DispensingStatus string
const (
	DispensingStatusBegin    DispensingStatus = "begin"
	DispensingStatusFinished                  = "finished"
	DispensingStatusFailed                    = "failed"
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
	Interval         int       `db:"interval"`
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
	db, err := sqlx.Connect(
		"postgres",
		"host=127.0.0.1 port=5432 user=box password=box dbname=box sslmode=disable",
	)
	return db, err
}
